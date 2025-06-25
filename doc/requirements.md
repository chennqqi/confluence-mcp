将原来github.com/ctreminiom/go-atlassian/confluence迁移到github.com/ctreminiom/go-atlassian/confluence/v2 

调用mcp在confluence上搜索icap 

在代码中增加调试信息代码
1. 定义环境变量ATLASSIAN_DEBUG，从此环境变量中获取DEBUG日志目录，默认值为I:\github.com\confluence-mcp\debug.txt
2. 自定义http.Client，将请求和响应输出输出到DEBUG文件中 

调用Confluence mcp搜索包含cursor的页面 