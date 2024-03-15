using System.Collections.Concurrent;

namespace Phoenix.Network
{
    // 缓存，等外部处理
    public class QueueProcessor : IMsgProcessor
    {
        //Queue<IMsg> _msgs = new Queue<IMsg>();  
        ConcurrentQueue<IMsg> _msgs = new ConcurrentQueue<IMsg>();  

        public void Process(IProtocol protocol, IMsg msg)
        {
            //lock(_msgs)
            //{
            //    _msgs.Enqueue(msg);
            //}
            _msgs.Enqueue(msg);
        }

        public void Update()
        {

        }

        public IMsg Pop()
        {
            //if (_msgs.Count == 0)
            //    return null;
            //lock(_msgs)
            //{
            //    try 
            //    {
            //        return _msgs.Dequeue();
            //    }
            //    catch(Exception e)
            //    {
            //        Console.WriteLine("Exception:\n" + e);
            //        return null;
            //    }                
            //}

            IMsg msg;
            if (!_msgs.TryDequeue(out msg))
                return null;
            return msg;
        }
    }

    public class QueueProcessorFactory : IProcessorFactory
    {
        public IMsgProcessor Create()
        {
            return new QueueProcessor();
        }
    }
}

