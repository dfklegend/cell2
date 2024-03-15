using UnityEngine;
using System;
using Phoenix.Core;
using UnityEngine.SceneManagement;

namespace Phoenix.Scene
{
    public class SceneCtrl : Singleton<SceneCtrl>
    {
        AsyncOperation _curOp;
        public void LoadSceneAsync(string scene, Action complete)
        {
            _curOp = SceneManager.LoadSceneAsync(scene);
            _curOp.completed += (op) => 
            {
                complete?.Invoke();
                _curOp = null;
            };
        }

        public bool IsLoading() 
        {
            return _curOp != null;
        }

        public float GetProgress()
        {
            if (_curOp == null)
                return 1f;
            return _curOp.progress;
        }

        public bool IsLoadOver()
        {
            return _curOp == null;
        }
    }
}

