需求分析：
用户希望通过mcp工具在Confluence上搜索关键词"icap"。
实现思路：
1. 调用mcp_confluence_search_page接口，传入关键词"icap"。
2. 返回搜索结果。
3. 结果可用于后续页面获取、评论等操作。

1. 检查所有引用github.com/ctreminiom/go-atlassian/confluence的地方。
2. 替换为github.com/ctreminiom/go-atlassian/confluence/v2。
3. 检查API是否有不兼容变更，必要时调整调用方式。

问题分析：
用户用curl手动请求Confluence搜索接口时有权限，但通过本项目mcp调用时返回401错误。
初步分析可能原因如下：
1. mcp项目代码中请求Confluence时未正确携带认证信息（如Authorization Header、Cookie等）。
2. mcp项目使用的认证方式与curl手动测试不一致（如curl用的是Basic Auth或Token，mcp用的是Session/Cookie等）。
3. mcp项目代码中httpclient实现存在bug，导致Header丢失或格式不正确。
4. 认证信息在mcp配置文件或环境变量中未正确传递到请求层。
5. Confluence对User-Agent、Referer等Header有限制，mcp请求被拦截。
6. 认证信息被中间件（如代理、网关）拦截或替换。
建议：
- 对比curl和mcp请求的所有Header，重点关注Authorization、Cookie、User-Agent等。
- 检查mcp项目中httpclient的实现，确认Header传递逻辑。
- 检查mcp配置文件/环境变量，确认认证信息来源。
- 可在mcp中增加请求日志，输出实际请求Header和响应内容。

需求分析：
1. 需要在项目中增加调试日志功能，便于排查Confluence API请求问题。
2. 通过环境变量ATLASSIAN_DEBUG指定日志文件路径，若未设置则使用默认路径I:\github.com\confluence-mcp\debug.txt。
3. 需自定义http.Client，在每次请求和响应时，将详细内容（如Header、Body、URL、状态码等）输出到日志文件。
4. 这样可方便对比curl和mcp请求的差异，定位401等问题。
实现建议：
- 封装一个自定义Transport，拦截RoundTrip方法，记录请求和响应。
- 日志内容建议包含时间、请求方法、URL、Header、Body、响应状态码、响应Header、响应Body等。
- 日志写入需考虑并发安全。

需求分析：
1. 用户希望通过Confluence MCP接口，搜索所有包含"cursor"关键字的页面。
2. 需要调用search-page相关工具或API，传入关键词"cursor"。
3. 返回结果应为包含该关键词的页面列表。
4. 需保证敏感信息不被读取。 