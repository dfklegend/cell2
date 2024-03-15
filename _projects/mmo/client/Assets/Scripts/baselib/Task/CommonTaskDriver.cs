using UnityEngine;
using System.Collections;
using System.Collections.Generic;



namespace Phoenix.AppTask
{    
    public class CommonTaskDriver : MonoBehaviour
    {
        public static CommonTaskDriver It = null;
        
        List<BaseCommonTask> _tasks = new List<BaseCommonTask>();
        // Use this for initialization
        void Awake()
        {
            It = this;            
        }     

        void Update()
        {   
            taskUpdate();            
        }

        void LateUpdate()
        {   
            taskLateUpdate();
        }
        
        public void AddTask(BaseCommonTask task)
        {
            _tasks.Add(task);
        }

        private void taskUpdate()
        {
            if (_tasks.Count == 0)
                return;
            for (int i = 0; i < _tasks.Count; i++)
                _tasks[i].Update();
        }

        private void taskLateUpdate()
        {
            if (_tasks.Count == 0)
                return;
            for (int i = 0; i < _tasks.Count; i++)
                _tasks[i].LateUpdate();
        }
    }
}
