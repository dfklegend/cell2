using System;

namespace Phoenix.Core
{
    // 非线程安全
    // 逻辑线程唯一分配一个
    /**
     * Delegate type for handle elapsed timer callback.
     */
    public delegate void TimerHandler( object[] ps );

    /**
     * @class TimerID
     *
     * @brief Wrapper class for timer id.
     */
    public class TimerID
    {
        /// Constants for invalid value.
        public static readonly TimerID INVALID_VALUE = new TimerID(Int64.MaxValue);

        /// Initializes.
        public TimerID(long identifier)
        {
            this.identifier_ = identifier;
        }

        public static implicit operator long(TimerID timerID)
        {
            return timerID.identifier_;
        }

        public static bool operator ==(TimerID lhs, TimerID rhs)
        {
            // If both are null, or both are same instance, return true.
            if (System.Object.ReferenceEquals(lhs, rhs))
            {
                return true;
            }

            // If one is null, but not both, return false.
            if (((object)lhs == null) || ((object)rhs == null))
            {
                return false;
            }

            // Return true if the fields match:
            return lhs.identifier_ == rhs.identifier_;
        }

        public static bool operator !=(TimerID lhs, TimerID rhs)
        {
            return !(lhs == rhs);
        }

        public override bool Equals(System.Object rhs)
        {
            // If parameter is null return false
            if (rhs == null)
            {
                return false;
            }

            // If parameter cannot be cast to TimerID return false
            TimerID timerID = rhs as TimerID;
            if ((System.Object)timerID == null)
            {
                return false;
            }

            // Returns true if the fields match
            return this.identifier_ == timerID.identifier_;
        }

        public bool Equals(TimerID rhs)
        {
            // If parameter is null return false
            if ((object)rhs == null)
            {
                return false;
            }

            // Return true if the fields match 
            return this.identifier_ == rhs.identifier_;
        }

        public override int GetHashCode()
        {
            return this.identifier_.GetHashCode();
        }

        /// Keep track of identifier for registered timer.
        private long identifier_;
    }

    /**
     * @class TimerNodeDispatchInfo
     *
     * @brief Maintains generated dispatch information for Timer nodes.
     */
    internal class TimerNodeDispatchInfo
    {
        /// The handler of object held in the queue.
        public TimerHandler handler_;

        /// Flag to check if the timer is recurring.
        public bool isRecurringTimer_;
        public object[] params_;

        public void Clear()
        {
            handler_ = null;
            params_ = null;
        }
    }

    /**
     * @class TimerNode
     *
     * @brief Maintains the state associated with a Timer entry.
     */
    internal class TimerNode
    {
        /// Initializes.
        public TimerNode(TimerHandler handler, TimerID timerID,
            TimerNode prev, TimerNode next, long future, float interval)
        {
            this.handler_ = handler;
            this.future_ = new DateTime(future);
            this.interval_ = interval;
            this.timerID_ = timerID;
            this.prev_ = prev;
            this.next_ = next;
        }

        /// 获取用于计时器到期处理所需要的相关信息。
        public void GetDispatchInfo(ref TimerNodeDispatchInfo info)
        {
            // Yes, do a copy.
            info.handler_ = this.handler_;
            info.isRecurringTimer_ = this.interval_ > 0.0f;
            info.params_ = _params;
        }

        #region Property Accessors

        public TimerHandler handler
        {
            get { return this.handler_; }
        }

        public DateTime future
        {
            get { return this.future_; }
            set { this.future_ = value; }
        }

        public float interval
        {
            get { return this.interval_; }
        }

        public TimerID timerID
        {
            get { return this.timerID_; }
        }

        public TimerNode Prev
        {
            get { return this.prev_; }
            set { this.prev_ = value; }
        }

        public TimerNode Next
        {
            get { return this.next_; }
            set { this.next_ = value; }
        }

        #endregion

        /// Handler that was stored in the timer queue.
        private TimerHandler handler_;

        /// Time until the timer expires.
        private DateTime future_;

        /// If this is a periodic timer this holds the time 
        /// until the next timeout.
        private float interval_;

        /// Id of this timer (used to cancel timers before they expire).
        private TimerID timerID_;

        /// Pointer to previous timer.
        private TimerNode prev_;

        /// Pointer to next timer.
        private TimerNode next_;

        private object[] _params;
        public void SetParams( object[] p )
        {
            _params = p;
        }

        public object[] GetParams()
        {
            return _params;
        }


    }

    /**
     * @class TimerManager
     *
     * @brief Provides a based list timer queue. 
     * It uses a list of absolute times. In the average case, scheduling and 
     * canceling timers is O(N) (where N is the total number of timers) and 
     * expiring timers is O(K) ( where K is the total number of timers that 
     * are < the current time of day).
     * 
     * 每帧一个timter只会跳一次(Expire 触发的)
     * 会根据时间差异，尽量补时间，保证下一次跳的间隔(但是只补一帧)
     * (recomputeNextAbsIntervalTime)
     */
    public class TimerManager
    {
        public static readonly DateTime ZERO = new DateTime(0);

        /// 计时器分辨率：0.01秒。
        public const float TIMER_RESOLUTION = 0.001f;

        /// <summary>
        /// 游戏时间
        /// </summary>
        public int totalGameTime = 0;

        private TimerNodeDispatchInfo _tmpDispInfo = new TimerNodeDispatchInfo();

        /// Close timer queue and cancels all timers.
        /// Returns number of timers cancelled.
        public int Close()
        {
            int numberOfTimersCancelled = 0;

            // Remove all remaining items in the list.
            while (this.removeFirst() != null)
            {
                ++numberOfTimersCancelled;
            }

            return numberOfTimersCancelled;
        }

        /// True if queue is empty, else false.
        public bool IsEmpty()
        {
            return this.head_ == this.head_.Next;
        }

        /**
         * 计时器注册。
         *
         * @param handler 计时器到期回调对象。
         * 
         * @param futureTime 注册到第一次回调之间的间隔时间。分辨率为 0.1 秒。
         * @param interval 周期性回调的间隔时间。
         *
         * @return 返回一个用于标识被添加的计时器的 TimerID 对象。
         * 该对象可以用在 cancel 调用来取消关联的计时器。
         */
        public TimerID AddTimer(TimerHandler handler, float futureTime, float interval, params object[] args)
        {
            if (handler == null)
            {
                PConsole.Log("ERROR: TimerManager::addTimer: handler is null!");

                return TimerID.INVALID_VALUE;
            }

            // 收拢时间间隔到最小分辨率
            if (futureTime > 0.0f && futureTime < TimerManager.TIMER_RESOLUTION)
            {
                futureTime = TimerManager.TIMER_RESOLUTION;
            }

            // 收拢时间间隔到最小分辨率
            if (interval > 0.0f && interval < TimerManager.TIMER_RESOLUTION)
            {
                interval = TimerManager.TIMER_RESOLUTION;
            }

            if (this.idCounter_ == Int64.MaxValue)
            {
                PConsole.Log("ERROR: TimerManager::addTimer: " + 
                    "id counter has out of range in TimerManager!");

                return TimerID.INVALID_VALUE;
            }

            TimerID timerID = new TimerID(this.idCounter_++);
            long tickFurture = DateTime.Now.Ticks + (long)(futureTime * 10000000);

            TimerNode node = new TimerNode(handler,
                timerID, null, null, tickFurture, interval);
            node.SetParams(args);

            this.schedule(node, node.future);

            return timerID;
        }

        /**
         * Callback the <TimerHandler> for all timers whose values are <= DateTime.Now.  
         *
         * Depending on the resolution of the underlying the system calls like 
         * might return at time different than that is specified in the timeout. 
         * Suppose the system guarantees a resolution of t ms.
         * The time line will look like
         *
         *             A                   B
         *             |                   |
         *             V                   V
         *  |-------------|-------------|-------------|-------------|
         *  t             t             t             t             t
         *
         *
         * If you specify a timeout value of A, then the timeout will not occur
         * at A but at the next interval of the timer, which is later than that 
         * is expected. Similarly, if your timeout value is equal to B, then the 
         * timeout will occur at interval after B. 
         *
         * Things get interesting if the t before the timeout value B is zero
         * i.e your timeout is less than the interval. In that case, you are
         * almost sure of not getting the desired timeout behaviour. Maybe you
         * should look for a better OS :-).
         *
         * Returns the number of timers expired.
         */
        public int Expire()
        {
            // Keep looping while there are timers remaining and the earliest
            // timer is <= the <currTime> passed in to the method.

            if (this.IsEmpty())
            {
                return 0;
            }

            int numberOfTimersExpired = 0;

            DateTime currTime = DateTime.Now;
            TimerNodeDispatchInfo info = _tmpDispInfo;

            while (this.dispatch(currTime, ref info))
            {
                TimerManager.upcall(info);

                ++numberOfTimersExpired;
            }
            _tmpDispInfo.Clear();

            return numberOfTimersExpired;
        }

        /**
         * Cancel the single timer that matches the @a <timerID> value (which
         * was returned from the <addTimer> method). Returns true if cancellation 
         * succeeded and false if the @a <timerID> wasn't found.
         */
        public bool Cancel(TimerID timerID)
        {
            TimerNode n = this.findNode(timerID);
            if (n == null)
            {
                return false;
            }

            this.unlink(ref n);

            return true;
        }

        /// Default constructor.
        public TimerManager()
        {
            this.head_ = new TimerNode(null, new TimerID(0), null, null, 0, 0);

            this.head_.Prev = this.head_;
            this.head_.Next = this.head_;

            this.idCounter_ = 1;
        }

        #region Implementations

        /// Reads the earliest node from the queue and returns it.
        private TimerNode first()
        {
            TimerNode n = this.head_.Next;
            if (n != this.head_)
            {
                return n;
            }
            return null;
        }

        /// Removes the earliest node from the queue and returns it.
        private TimerNode removeFirst()
        {
            TimerNode n = this.first();
            if (n != null)
            {
                this.unlink(ref n);
            }
            return n;
        }

        private void schedule(TimerNode n, DateTime expire)
        {
            if (this.IsEmpty())
            {
                n.Prev = this.head_;
                n.Next = this.head_;
                this.head_.Prev = n;
                this.head_.Next = n;
            }
            else
            {
                // We always want to search backwards from the tail of the list, because
                // this minimizes the search in the extreme case when lots of timers are
                // scheduled for exactly the same time, and it also assumes that most of
                // the timers will be scheduled later than existing timers.
                TimerNode curr = this.head_.Prev;
                while (curr != this.head_ && curr.future > expire)
                {
                    curr = curr.Prev;
                }

                // Insert after.
                n.Prev = curr;
                n.Next = curr.Next;
                curr.Next.Prev = n;
                curr.Next = n;
            }
        }

        /**
         * Get the dispatch information for a timer whose value is <= @a <currTime>.
         * Returns true if there is a node whose value <= @a <currTime> else returns a false.
         */
        private bool dispatch(DateTime currTime, ref TimerNodeDispatchInfo info)
        {
            if (this.IsEmpty())
            {
                return false;
            }

            if (this.earliestTime() <= currTime)
            {
                // Get the first timer node.
                TimerNode expired = this.removeFirst();

                // Get the dispatch info.
                expired.GetDispatchInfo(ref info);

                // Check if this is an interval timer.
                if (info.isRecurringTimer_)
                {
                    // Make sure that we skip past values that have already "expired".
                    this.recomputeNextAbsIntervalTime(ref expired, currTime);

                    // Since this is an interval timer, we need to reschedule it.
                    this.schedule(expired, expired.future);
                }

                return true;
            }

            return false;
        }

        /// Returns the time of the earlier node in the TimerManager.
        /// Must be called on a non-empty queue.
        private DateTime earliestTime()
        {
            TimerNode n = this.first();
            if (n != null)
            {
                return n.future;
            }
            return TimerManager.ZERO;
        }

        /// Recompute when the next time is that this interval timer should fire.
        private void recomputeNextAbsIntervalTime(ref TimerNode expired, DateTime currTime)
        {
            if (expired.future <= currTime)
            {
                // Compute the span between the current time and when the timer 
                // would have expired in the past (and normalize to nanoseconds).
                // 由于 .NET 下 DateTime 的 Ticks 的分辨率为 100-nanosecond，
                // 所以这里的公式为:
                //                           msec   usec   nsec
                //   nanoseconds = seconds * 1000 * 1000 * 10.
                long intervalNsecs = ((long)(expired.interval * 1000.0f)) * 10000L;
                TimeSpan diff = currTime - expired.future;

                // Compute the delta time in the future when the timer
                // should fire as if it had advanced incrementally.
                // The modulo arithmetic accomodates the likely case that
                // the current time doesn't fall precisely on a timer 
                // firing interval.
                // 最多追一帧
                long newIntervalNsecs = intervalNsecs - (diff.Ticks % intervalNsecs);

                // Compute the absolute time in the future when this interval timer 
                // should expire.
                expired.future = new DateTime(currTime.Ticks + newIntervalNsecs);
            }
        }

        /// 获得 @a <timerID> 关联的计时器节点。
        /// 如果 @a <timerID> 对应的节点不存在返回 null。
        private TimerNode findNode(TimerID timerID)
        {
            TimerNode n = this.first();
            if (n == null)
            {
                return null;
            }

            for (; n != this.head_; n = n.Next)
            {
                if (n.timerID == timerID)
                {
                    return n;
                }
            }

            return null;
        }

        /// 断开节点的前后关联，并释放相应的句柄引用计数。
        private void unlink(ref TimerNode n)
        {
            n.Prev.Next = n.Next;
            n.Next.Prev = n.Prev;
            n.Prev = null;
            n.Next = null;
        }

        /// This method will call the TimerHandler() object.
        private static void upcall(TimerNodeDispatchInfo info)
        {
            try
            {
                info.handler_(info.params_);
            }
            catch (Exception ex)
            {
                PConsole.Log("TimerManager::upcall throw an exception: "
                    + ex.Message + "\n" + ex.StackTrace);
            }
        }

        #endregion

        /// Pointer to linked list of <TimerHandles>.
        private TimerNode head_;

        /**
         * Keeps track of the timer id that uniquely identifies each timer.
         * This id can be used to cancel a timer via the <cancel(long)>
         * method.
         */
        private long idCounter_;
    }
}