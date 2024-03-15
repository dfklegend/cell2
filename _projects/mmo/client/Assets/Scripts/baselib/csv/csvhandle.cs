using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using UnityEngine;

namespace Phoenix.csv
{
    /// <summary>
    /// Class to store one CSV row
    /// </summary>
    public class CsvRow : List<string>
    {
        public string LineText { get; set; }
    }

    /// <summary>
    /// Class to write data to a CSV file
    /// </summary>
    public class CsvFileWriter : StreamWriter
    {
        public CsvFileWriter(Stream stream)
            : base(stream)
        {
        }

        public CsvFileWriter(string filename)
            : base(filename)
        {
        }

        /// <summary>
        /// Writes a single row to a CSV file.
        /// </summary>
        /// <param name="row">The row to be written</param>
        public void WriteRow(CsvRow row)
        {
            StringBuilder builder = new StringBuilder();
            bool firstColumn = true;
            foreach (string value in row)
            {
                // Add separator if this isn't the first value
                if (!firstColumn)
                    builder.Append(',');
                // Implement special handling for values that contain comma or quote
                // Enclose in quotes and double up any double quotes
                if (value.IndexOfAny(new char[] { '"', ',' }) != -1)
                    builder.AppendFormat("\"{0}\"", value.Replace("\"", "\"\""));
                else
                    builder.Append(value);
                firstColumn = false;
            }
            row.LineText = builder.ToString();
            WriteLine(row.LineText);
        }
    }    

    // 
    public class CSVReader
    {
        public class Row
        {
            public List<string> cols = new List<string>();
        }
        public List<Row> lines = new List<Row>();
        public void ParseFile(string strContent)
        {            
            int pos = 0;
            int iQuoteCount = 0;
            int nRowBegin = -1;
            int nRowEnd = -1;
            int iRowIndex = 0;

            nRowBegin = 0;
            while (pos < strContent.Length)
            {
                if (strContent[pos] == '"')
                    iQuoteCount++;
                else
                    if (strContent[pos] == '\n')
                    {
                        if (iQuoteCount % 2 == 0)
                        {
                            // 新的一行
                            nRowEnd = pos - 1;

                            OnGetRow(iRowIndex++, strContent, nRowBegin, nRowEnd);
                            nRowBegin = pos + 1;
                            nRowEnd = -1;
                        }
                    }
                pos++;
            }
            OnGetRow(iRowIndex++, strContent, nRowBegin, pos - 1);
        }

        void OnGetRow(int iRowIndex, string strContent, int nRowBegin, int nRowEnd)
        {
            if (nRowEnd < nRowBegin)
                return;
            // 去掉一个\r
            while (nRowEnd > nRowBegin)
            {
                if (strContent[nRowEnd] == '\r')
                    nRowEnd--;
                else
                    break;
            }
            string strRow = strContent.Substring(nRowBegin, nRowEnd - nRowBegin + 1);
            ParseRow(iRowIndex, strRow);
        }

        void ParseRow(int iRowIndex, string strRow)
        {
            int pos = 0;
            int rows = 0;
            bool bHasLastCol = false;
            List<string> cols = new List<string>();

            while (pos < strRow.Length)
            {
                string value;

                // Special handling for quoted field
                if (strRow[pos] == '"')
                {
                    // Skip initial quote
                    pos++;

                    // Parse quoted value
                    int start = pos;
                    while (pos < strRow.Length)
                    {
                        // Test for quote character
                        if (strRow[pos] == '"')
                        {
                            // Found one
                            pos++;

                            // If two quotes together, keep one
                            // Otherwise, indicates end of value
                            if (pos >= strRow.Length || strRow[pos] != '"')
                            {
                                pos--;
                                break;
                            }
                        }
                        pos++;
                    }
                    value = strRow.Substring(start, pos - start);
                    value = value.Replace("\"\"", "\"");
                }
                else
                {
                    // Parse unquoted value
                    int start = pos;
                    while (pos < strRow.Length && strRow[pos] != ',')
                        pos++;
                    value = strRow.Substring(start, pos - start);
                }

                // Add field to list
                cols.Add(value);
                rows++;

                // Eat up to and including next comma
                while (pos < strRow.Length && strRow[pos] != ',')
                    pos++;

                if (pos < strRow.Length && strRow[pos] == ',')
                    bHasLastCol = true;
                else
                    bHasLastCol = false;
                // skip ,
                if (pos < strRow.Length)
                    pos++;
            }// where

            if (bHasLastCol)
            {
                cols.Add(string.Empty);
                rows++;
            }
            // Delete any unused items
            while (cols.Count > rows)
                cols.RemoveAt(rows);

            // Return true if any columns read
            if (cols.Count > 0)
            {   
                OnRowParsed(iRowIndex, cols);                
            }
        }

        void OnRowParsed( int iIndex,List<string> cols )
        {   
            /*int i = 0 ;
            for (; i < cols.Count; i++)
            {
                if (cols[i].Length != 0)
                    break;
            }
            if (i == cols.Count)
                return;*/
            Row r = new Row();
            r.cols = cols;
            lines.Add(r);
        }
    }    
}