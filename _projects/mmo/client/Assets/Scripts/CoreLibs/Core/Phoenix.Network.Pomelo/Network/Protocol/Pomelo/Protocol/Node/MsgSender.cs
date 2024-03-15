using System;
using System.Collections.Generic;
using System.Threading;
using Pomelo.DotNetClient;
using TimeUtil = Phoenix.Utils.TimeUtil;

namespace Phoenix.Network.Protocol.Pomelo
{
    public enum RequestErrCode
    {
        Succ = 0,       // 成功
        Timeout,        // 超时        
        NotFoundOrBadRet,  // 不存在或者无效返回
        NotConnectted,  // 没连接
        Canceled,       // 被取消(连接中断了)
        NoMailBox,      // 没有目标邮箱
        InvalidRet,     // 无效返回
    }
    public class RequestResult
    {
        public uint reqId = 0;
        public RequestErrCode code;
        public byte[] rawData;
        // 如果需要反序列化成对象
        // 可以定制反序列化器
        // Succ才有效
        public object data;

        public bool IsSucc()
        {
            return code == RequestErrCode.Succ;
        }

        public T GetData<T>()
            where T: class
        {
            if (data == null)
                return null;
            return data as T;
        }

        public static RequestResult MakeErrResult(RequestErrCode code)
        {
            var result = new RequestResult();
            result.code = code;
            return result;
        }
    }

    public interface IClientSession
    {
        int GetHandle();
        bool IsConnected();
        void SendRequest(string route, uint reqId, byte[] data);
        void SendResponse(uint reqId, byte[] result);
        void SendNotify(string route, byte[] data);
    }

    public class RequestWaitAck
    {
        public float timeTimeout;
        public uint reqId;
        public Action<RequestResult> cb;

        // 结果类型
        public Type resultType;
    }

    // 负责发送Request，处理Response
    public class MsgSender
    { 
        private const float TIME_TIMEOUT = 15f;        
        //int _nextId = 1;

        private Dictionary<uint, RequestWaitAck> _callBackMap = new Dictionary<uint, RequestWaitAck>();
        private float _nextCheckExpired = 0;
        private Serializer.ISerializer _serializer;

        public void SetSerializer(Serializer.ISerializer serializer)
        {
            _serializer = serializer;
        }

        private uint allocReqId()
        {
            return gAllocReqId();
            //var id = _nextId++;
            //if (_nextId > 0xFFFFFFF)
            //    _nextId = 1;
            //return (uint)id;
        }

        private static int _gNextId = 1;
        private static uint gAllocReqId()
        {
            var id = Interlocked.Increment(ref _gNextId);
            if (_gNextId > 0xFFFFFFF)
                _gNextId = 1;
            return (uint)id;
        }

        public static void ResetReqId()
        {
            _gNextId = 1;
        }

        private RequestWaitAck doRequest(IClientSession session, string route, byte[] data, Action<RequestResult> cb)
        {
            if (data == null)
                data = PomeloDefine.empteBytes;
            RequestWaitAck req = new RequestWaitAck();
            // 超时
            req.timeTimeout = TimeUtil.GetSystemSecond() + TIME_TIMEOUT;
            req.reqId = allocReqId();
            req.cb = cb;

            _callBackMap[req.reqId] = req;
            // do send
            if (session.IsConnected())
            {
                session.SendRequest(route, req.reqId, data);
                //Env.L.FileLog($"{session.GetHandle()} reqId: {req.reqId} send");
            }
            else 
            {
                //Env.L.FileLog($"reqId: {req.reqId} not send");
            }
            

            return req;
        }

        public void Request(IClientSession session, string route, byte[] data, Action<RequestResult> cb)
        {
            doRequest(session, route, data, cb);        
        }

        public void Request(IClientSession session, string route, object arg, Action<RequestResult> cb)
        {
            if (_serializer == null)
                return;
            var data = _serializer.Serialize(arg);
            doRequest(session, route, data, cb);            
        }

        public void Request(IClientSession session, string route, object arg, Type resultType, Action<RequestResult> cb)
        {
            if (_serializer == null)
                return;
            var data = _serializer.Serialize(arg);
            var req = doRequest(session, route, data, cb);
            // 记录结果类型
            req.resultType = resultType;
        }

        public void Request<TRESULT>(IClientSession session, string route, object arg, Action<RequestResult> cb)
        {
            Request(session, route, arg, typeof(TRESULT), cb);
        }

        public void Notify(IClientSession session, string route, byte[] data)
        {
            if (data == null)
                data = PomeloDefine.empteBytes;
            if (session.IsConnected())
                session.SendNotify(route, data);
        }

        public void Notify(IClientSession session, string route, object arg)
        {
            if (_serializer == null)
                return;
            var data = _serializer.Serialize(arg);
            Notify(session, route, data);
        }

        public void ProcessResponse(Message response)
        {
            if(response.type != MessageType.MSG_RESPONSE)
            {
                throw new Exception("Must be MSG_RESPONSE");
            }
            RequestWaitAck req;
            if (!_callBackMap.TryGetValue(response.id, out req))
                return;
            RequestResult result = new RequestResult();
            result.code = RequestErrCode.Succ;
            if(response.error)
                result.code = RequestErrCode.NotFoundOrBadRet;
            result.rawData = response.rawData;
            result.reqId = response.id;

            // 序列化成目标结构
            if(_serializer != null && req.resultType != null)
            {
                // 一般是找不到处理者
                // protobuf如果数据都是缺省值，会序列化出空数据                
                result.data = _serializer.Deserialize(result.rawData, req.resultType);
            }
            _callBackMap.Remove(response.id);
            
            req.cb(result);            
        }

        public void Update()
        {
            // 处理超时
            processRequestTimeout();
        }

        private void processRequestTimeout()
        {
            var now = TimeUtil.GetSystemSecond();            
            if (now < _nextCheckExpired)
                return;
            _nextCheckExpired = now + 5f;

            List<uint> toRemove = new List<uint>();
            foreach (KeyValuePair<uint, RequestWaitAck> kv in _callBackMap)
            {
                var req = kv.Value;
                if(now > req.timeTimeout)
                {             
                    toRemove.Add(kv.Key);
                }
            }

            for (var i = 0; i < toRemove.Count; i++)
            {
                var reqId = toRemove[i];
                var req = _callBackMap[reqId];
                RequestResult result = new RequestResult();
                result.reqId = reqId;
                result.code = RequestErrCode.Timeout;
                req.cb(result);

                _callBackMap.Remove(reqId);
            }
        }

        public void ProcessRequestNotFound(uint reqId)
        {
            RequestWaitAck req;
            if (!_callBackMap.TryGetValue(reqId, out req))
                return;

            RequestResult result = new RequestResult();
            result.code = RequestErrCode.NotFoundOrBadRet;            
            result.reqId = reqId;            

            req.cb(result);
            _callBackMap.Remove(reqId);
        }

        public void Stop()
        {
            cancelAllRequestWait();
        }

        private void cancelAllRequestWait()
        {   
            foreach (KeyValuePair<uint, RequestWaitAck> kv in _callBackMap)
            {
                var req = kv.Value;

                RequestResult result = RequestResult.MakeErrResult(RequestErrCode.Canceled);                
                req.cb(result);
            }

            _callBackMap.Clear();
        }
    }
}

