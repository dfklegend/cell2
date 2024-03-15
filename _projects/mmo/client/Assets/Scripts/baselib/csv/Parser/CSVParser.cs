using Phoenix.Res;
using System;
using System.Collections.Generic;
using System.Reflection;
using UnityEngine;

namespace Phoenix.csv
{
    public static class CSVParser
    {
        static void ParsePropertyValue<T>(int row,int col,T obj, FieldInfo fieldInfo, string valueStr)
        {
            if (valueStr == null)
                valueStr = string.Empty;
            System.Object value = null;

            BaseValueConvert converter = CSVValueConverter.It.FindConvert(fieldInfo.FieldType);
            // 优先走converter
            if (fieldInfo.FieldType.IsEnum && converter == null)
            {
                try
                {
                    value = Enum.Parse(fieldInfo.FieldType, valueStr);
                }
                catch (System.Exception)
                {
                    Log.LogCenter.Asset.Debug("(row:{0},col:{1}) is not correct Enum", row, col);
                }
            }
            else
            {   
                if (converter != null)
                {
                    value = converter.Convert(valueStr);
                }
                else
                {                    
                    if (valueStr.Contains("\"\""))
                        valueStr = valueStr.Replace("\"\"", "\"");

                    // process the excel string.
                    if (valueStr.Length > 2 && valueStr[0] == '\"' && valueStr[valueStr.Length - 1] == '\"')
                        valueStr = valueStr.Substring(1, valueStr.Length - 2);
                    value = valueStr;                    
                }
            }

            if (value == null)
                return;

            fieldInfo.SetValue(obj, value);
        }

        static T ParseObject<T>(int iLineNo,string[] values, FieldInfo[] propertyInfos)
        {
            if (values.Length < propertyInfos.Length)
                return default(T);
            T obj = Activator.CreateInstance<T>();            
            for (int j = 0; j < propertyInfos.Length; j++)
            {
                if (propertyInfos[j] == null)
                    continue;
                string value = values[j];                
                try
                {
                    ParsePropertyValue(iLineNo,j+1,obj, propertyInfos[j], value);
                }
                catch (Exception ex)
                {
                    Debug.LogError("line:" + values + " for: " + propertyInfos[j].Name+" "+j);
                    Debug.LogError(ex);
                }
            }
            return obj;
        }

        static FieldInfo[] GetPropertyInfos<T>(string memberLine)
        {
            string[] members = memberLine.Split(",".ToCharArray(), StringSplitOptions.RemoveEmptyEntries);
            return GetPropertyInfos<T>(members);
        }

        static FieldInfo[] GetPropertyInfos<T>(string[] cols)
        {
            Type objType = typeof(T);
            
            FieldInfo[] propertyInfos = new FieldInfo[cols.Length];
            for (int i = 0; i < cols.Length; i++)
            {
                propertyInfos[i] = objType.GetField(cols[i]);
                if (propertyInfos[i] == null)
                {
                    //Debug.LogError("GetField is None!" + cols + " " + i + " " + cols[i]);
                }
            }

            return propertyInfos;
        }

        private static string loadTextFile(string path)
        {
            return ResourceMgr.It.LoadTextFile(path);
        }

        static public T[] Parse<T>(string path, int indexRead)
        {
            // we load the table data from package. from udata
            // here we load the text asset.            
            string strContent = loadTextFile(path);
            if (strContent == null)
            {
                Debug.LogError("无法加载表格文件：" + path);
                return null;
            }
            return ParseFromStr<T>(strContent, indexRead);
        }
        
        static public T[] ParseFromStr<T>(string strContent,int indexRead)
        {           
            CSVReader cn = new CSVReader();
            cn.ParseFile(strContent);
            List<CSVReader.Row> lines = cn.lines;
            if (lines.Count < 3)
            {
                //Debug.LogError("表格文件行数错误，【1】属性名称【2】变量名称【3-...】值：" + path);
                return null;
            }
            // fetch all of the field infos.
            FieldInfo[] propertyInfos = GetPropertyInfos<T>(lines[indexRead-1].cols.ToArray());
            // parse it one by one.

            List<T> objs = new List<T>();
            T obj;
            for (int i = 0; i < lines.Count - indexRead; i++)
            {
                obj = ParseObject<T>(i + indexRead + 1, lines[i + indexRead].cols.ToArray(), propertyInfos);
                if (obj == null)
                    continue;
                objs.Add(obj);
            }
            return objs.ToArray();
        }
    }
    
}