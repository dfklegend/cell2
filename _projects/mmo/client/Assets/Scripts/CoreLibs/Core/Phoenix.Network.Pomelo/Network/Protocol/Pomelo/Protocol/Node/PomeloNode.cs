using System;
using System.Threading;
using System.Threading.Tasks;
using Pomelo.DotNetClient;

namespace Phoenix.Network.Protocol.Pomelo
{
    public class PomeloNode
    {
        IClientSession _session;
        public IClientSession session { get { return _session; } }
        // 可以多个Session共用
        // Update
        IMsgDispatcher _dispatcher;
        // 每个Session一个
        MsgSender _sender;        

        public MsgSender sender { get { return _sender; } }

        public void Init(IClientSession session, MsgSender sender)
        {
            _session = session;
            _sender = sender;            
        }

        public void SetMsgDispatcher(IMsgDispatcher processor)
        {
            _dispatcher = processor;
        }        

        public void Process(Message msgIn)
        {
            //Env.L.FileLog($"{_session.GetHandle()} process type: {msgIn.type} reqId: {msgIn.id} {msgIn.route}");
            //Console.WriteLine($"{_session.GetHandle()} process type: {msgIn.type} reqId: {msgIn.id} {msgIn.route}");
            if (msgIn.type == MessageType.MSG_RESPONSE)
            {
                _sender.ProcessResponse(msgIn);
                return;
            }
            if (processSystemNotify(msgIn))
                return;
            _dispatcher.Process(_session, msgIn);
        }

        private T tryParse<T>(byte[] rawData)
            where T: class
        {
            var serializer = _dispatcher.GetSerializer();
            if (serializer == null)
                return null;
            return serializer.Deserialize(rawData, typeof(T)) as T;
        }

        private bool processSystemNotify(Message msgIn)
        {
            if (msgIn.route == "__req_notfound__")
            {                
                var obj = tryParse<RequestNotFound>(msgIn.rawData);
                if (obj == null)
                    return true;
                Console.WriteLine($"Request not found: {obj.info}");
                _sender.ProcessRequestNotFound(obj.reqId);
                return true;
            }

            if(msgIn.route == "__error__")
            {
                var obj = tryParse<string>(msgIn.rawData);
                if (obj == null)
                    return true;
                Console.WriteLine($"Request not found: {obj}");                
            }
            return false;
        }

        public void Request(string route, byte[] data, Action<RequestResult> cb)
        {
            _sender.Request(_session, route, data, cb);
        }

        public void Request(string route, object data, Action<RequestResult> cb)
        {
            _sender.Request(_session, route, data, cb);
        }

        public async Task<RequestResult> RequestAsyncTask(string route, byte[] data)
        {
            var task = new RequestTask(this, route, data);
            await task;
            return task.GetResult();
        }

        public Task<RequestResult> RequestAsync(string route, object data)
        {
            var tcs = new TaskCompletionSource<RequestResult>();
            _sender.Request(_session, route, data, (result)=> {
                //Console.WriteLine($"tcs.SetResult Thread:{Thread.CurrentThread.ManagedThreadId}");
                tcs.SetResult(result);
            });
            //tcs.Task.ConfigureAwait(false);
            return tcs.Task;
        }

        public void Notify(string route, byte[] data)
        {
            _sender.Notify(_session, route, data);
        }

        public void Notify(string route, object arg)
        {
            _sender.Notify(_session, route, arg);
        }

        public void Update()
        {
            _sender.Update();           
        }

        public void Start()
        {
            //_processor.OnSessionStart(_session);
        }

        public void Stop()
        {
            _sender.Stop();
            _dispatcher.OnSessionStop(_session);
        }       

        // 
        public void Request<TRESULT>(string route, object arg, Action<RequestResult> cb)
        {
            _sender.Request<TRESULT>(_session, route, arg, cb);
        }

        public void Request(string route, object arg, Type resultType, Action<RequestResult> cb)
        {
            _sender.Request(_session, route, arg, resultType, cb);
        }

        public Task<RequestResult> RequestAsync<TRESULT>(string route, object arg)
        {
            var tcs = new TaskCompletionSource<RequestResult>();
            _sender.Request<TRESULT>(_session, route, arg, (result) => {
                //Console.WriteLine($"tcs.SetResult Thread:{Thread.CurrentThread.ManagedThreadId}");
                // 注: TCS.SetResult同步调用
                try
                {
                    tcs.TrySetResult(result);
                }
                catch(Exception e)
                {
                    Console.WriteLine(e);
                }
                
            });
            //tcs.Task.ConfigureAwait(false);            
            return tcs.Task;
        }
    }
}

