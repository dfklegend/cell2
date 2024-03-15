using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Starter : MonoBehaviour
{
    ClientApp _app;
    // Start is called before the first frame update
    void Start()
    {
        ClientApp app = new ClientApp();                
        app.Prepare();

        Phoenix.Log.LogCenter.Default.Debug("Start");
        app.Start();
        _app = app;
    }

    // Update is called once per frame
    void Update()
    {
        _app.Update();        
    }

    private void OnApplicationQuit()
    {
        Debug.Log("OnApplicationQuit");
        _app.OnQuit();
    }
}
