using System;
using System.Collections.Generic;
using System.Text;
using Pomelo.DotNetClient;
using TimeUtil = Phoenix.Utils.TimeUtil;

namespace Phoenix.Network.Protocol.Pomelo
{
    using NotifyAction = Action<IClientSession, object>;
    using RequestAction = Action<IClientSession, object, Action<object>>;

    public class IncomeRequest
    {
        public uint incomeReqId;
        public int serialId;
        public string route;
        // 过期直接移除，60s
        public float timeExpire;
        public IClientSession session;
    }    

    // 消息派发
    // 处理来访的Request
    // 结果回来之后，转换为Response回调回去
    // 可以注册回调函数
    // 也可以使用APIService的方式来定义处理器
    public class RichMsgDispatcher : IMsgDispatcher
    {
        const float TIME_EXPIRE = 60f;        

        int _nextSerialId = 1;

        // 请求
        Dictionary<int, IncomeRequest> _requests =
            new Dictionary<int, IncomeRequest>();       

        private float _nextCheckExpired = 0;

        private Serializer.ISerializer _serializer;
        
        private FuncDispatcherImpl _func = new FuncDispatcherImpl();
        private IDispatcherImpl _service;

        private int allocSerialId()
        {
            var id = _nextSerialId++;
            if (_nextSerialId > 0xFFFFFFF)
                _nextSerialId = 1;
            return id;
        }

        public void CreateServiceImpl(string categories)
        {
            _service = new ServiceDispatcherImpl(categories);
            updateSerializer(_serializer);
        }

        public void SetServiceImpl(IDispatcherImpl service)
        {
            _service = service;
            updateSerializer(_serializer);
        }

        public void SetSerializer(Serializer.ISerializer serializer)
        {
            _serializer = serializer;
            updateSerializer(serializer);
        }

        public Serializer.ISerializer GetSerializer()
        {
            return _serializer;
        }

        private void updateSerializer(Serializer.ISerializer serializer)
        {
            _func.SetSerializer(serializer);
            _service?.SetSerializer(serializer);
        }        

        // hander 注册接口
        public void RegisterRequestHandler(string route, RequestAction action)
        {
            _func.RegisterRequestHandler(route, action);            
        }

        public void RegisterRequestHandler<T>(string route, RequestAction action)
        {
            _func.RegisterRequestHandler<T>(route, action);           
        }

        public void RegisterNotifyHandler(string route, NotifyAction action)
        {
            _func.RegisterNotifyHandler(route, action);            
        }

        public void RegisterNotifyHandler<T>(string route, NotifyAction action)
        {
            _func.RegisterNotifyHandler<T>(route, action);
        }

        public void Process(IClientSession session, Message msgIn)
        {
            if(msgIn.type == MessageType.MSG_RESPONSE)
            {
                throw new Exception("can not handle response");                
            }            

            if (msgIn.type == MessageType.MSG_REQUEST)
            {
                processRequest(session, msgIn);
                return;
            }

            // MSG_NOTIFY or MSG_PUSH
            processNotify(session, msgIn);
        }

        private IDispatcherImpl getImpl(string route)
        {
            if (_func.HasHandler(route))
            {
                return _func;
            }
            if (_service != null && _service.HasHandler(route))
                return _service;
            return null;
        }

        private void processRequest(IClientSession session, Message msgIn)
        {
            //RequestHandler handler;
            //_handlers.TryGetValue(msgIn.route, out handler);
            //if (handler == null)
            //{
            //    // 没有定义
            //    sendError(session, $"can not find handler:{msgIn.route}");
            //    return;
            //}
            var impl = getImpl(msgIn.route);

            if(null == impl)
            {
                // 没有定义
                //sendError(session, $"can not find handler:{msgIn.route}");
                sendRequestNotFound(session, msgIn.id, $"can not find handler:{msgIn.route}");
                return;
            }

            if(msgIn.id == 0)
            {
                return;
            }

            int serialId = allocSerialId();

            var req = new IncomeRequest();
            req.session = session;
            req.serialId = serialId;
            req.incomeReqId = msgIn.id;
            req.route = msgIn.route;
            req.timeExpire = TimeUtil.Now() + TIME_EXPIRE;

            _requests[req.serialId] = req;

            //Env.L.FileLog($"serialId: {serialId} reqId: {req.incomeReqId} call");

            impl.InvokeRequest(session, msgIn.route, msgIn.rawData, (result) => {
                processResult(req.serialId, result);
            });

            // 序列化成目标结构
            //object data = msgIn.rawData;
            //if (_serializer != null && handler.argType != null)
            //{
            //    msgIn.data = _serializer.Deserialize(msgIn.rawData, handler.argType);
            //    data = msgIn.data;
            //}

            //// 调用handler
            //handler.cb.Invoke(session, data, (result) =>
            //{
            //    processResult(req.serialId, result);
            //});
        }

        private void processNotify(IClientSession session, Message msgIn)
        {            
            var impl = getImpl(msgIn.route);
            if (null == impl)                
            {
                Env.L.Error($"can not find handler:{msgIn.route}");
                // 没有定义
                sendError(session, $"can not find handler:{msgIn.route}");
                return;
            }

            impl.InvokeNotify(session, msgIn.route, msgIn.rawData);
            return;           
        }        

        private void processResult(int serialId, object result)
        {
            IncomeRequest req;
            _requests.TryGetValue(serialId, out req);
            if (req == null)
            {
                // 超时或者连接中断了
                Console.WriteLine($"processResult, can not find incomeRequest of {serialId}");
                return;
            }

            if (result == null)
                result = PomeloDefine.empteBytes;
            //Env.L.FileLog($"serialId: {serialId} reqId: {req.incomeReqId} ret");
            //
            byte[] body = null;
            if (result.GetType() == typeof(byte[]))
                body = (byte[])result;
            else
            if (_serializer != null)
                body = _serializer.Serialize(result);            
          
            sendResponse(req.session, req.incomeReqId, body);
            _requests.Remove(serialId);
        }

        private void sendResponse(IClientSession session, uint reqId, byte[] result)
        {
            if (!session.IsConnected())
                return;
            if (result == null)
                result = PomeloDefine.empteBytes;
            session.SendResponse(reqId, result);
        }

        private void sendError(IClientSession session, string error)
        {
            if (!session.IsConnected())
                return;
            if(_serializer != null)
            {
                // 序列化
                //session.SendNotify("__error__", _serializer.Serialize(error));
                return;
            }
            //session.SendNotify("__error__", Encoding.UTF8.GetBytes(error));
        }

        private void sendRequestNotFound(IClientSession session, uint reqId, string error)
        {
            if (!session.IsConnected())
                return;
            if (_serializer != null)
            {
                var msg = new RequestNotFound() { reqId = reqId, info = error };
                session.SendNotify("__req_notfound__", _serializer.Serialize(msg));
            }
        }

        // 移除过期的IncomeRequest
        public void Update()
        {
            removeExpiredIncomeReq();
        }

        private void removeReqs(Func<IncomeRequest, object, bool> func, object arg)
        {
            List<int> toRemove = new List<int>();
            foreach (KeyValuePair<int, IncomeRequest> kv in _requests)
            {
                var one = kv.Value;
                if (func(one, arg) )
                {
                    toRemove.Add(kv.Key);
                }
            }

            for (var i = 0; i < toRemove.Count; i++)
            {
                _requests.Remove(toRemove[i]);
            }
        }

        private void removeExpiredIncomeReq()
        {
            var now = TimeUtil.GetSystemSecond();
            if (now < _nextCheckExpired)
                return;
            _nextCheckExpired = now + 5f;
            
            removeReqs((req, arg) =>
            {
                if (req.timeExpire < now)
                {
                    Console.WriteLine($"req {req.route} expired");
                    return true;
                }
                return false;
            }, null);
        }

        public void OnSessionStart(IClientSession session)
        {
            //
        }

        // 清理所属session的请求
        public void OnSessionStop(IClientSession session)
        {
            removeReqs((req, arg) => {
                if (req.session == session)
                    return true;
                return false;
            }, session);
        }
    }
}

