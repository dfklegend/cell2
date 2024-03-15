using UnityEngine;
using System.Collections;

// 使节点保持不删除
public class DontDestoryed : MonoBehaviour 
{
	// Use this for initialization
	void Start () 
    {
        DontDestroyOnLoad(gameObject);
		Destroy(this);
	}
	
	// Update is called once per frame
	void Update () 
    {
	
	}
}
