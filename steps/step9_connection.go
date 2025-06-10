package steps

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"macos-clodop-schoolpal/config"
)

// runCommand 执行系统命令
func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// TestConnection 测试打印机连接
func TestConnection(cfg *config.Config) error {
	localPort := cfg.Network.LocalPort
	remoteHost := cfg.Network.RemoteHost
	remotePort := cfg.Network.RemotePort

	fmt.Println("🔗 测试网络连接...")

	// 测试本地端口转发是否正常
	err := testLocalPort(localPort)
	if err != nil {
		return fmt.Errorf("本地端口测试失败: %v", err)
	}
	fmt.Println("✅ 本地端口连接正常")

	// 测试远程连接是否可达
	err = testRemoteConnection(remoteHost, remotePort)
	if err != nil {
		return fmt.Errorf("远程连接测试失败: %v", err)
	}
	fmt.Println("✅ 远程连接正常")

	// 智能检测Clodop服务是否可用
	fmt.Println("🖨️ 检测Clodop服务...")

	// 将字符串端口转换为整数
	portInt := 0
	if localPort != "" {
		if p, err := strconv.Atoi(localPort); err == nil {
			portInt = p
		}
	}

	clodopPort, err := detectClodopPort(portInt)
	if err != nil {
		fmt.Printf("⚠️ Clodop服务检测失败: %v\n", err)
		fmt.Println("💡 这可能是因为:")
		fmt.Println("   - 远程Windows电脑上Clodop服务未运行")
		fmt.Println("   - 打印机未连接或未开机")
		fmt.Println("   - VPN连接不稳定")
		fmt.Println("   - 端口转发配置有问题")
		// 不要返回错误，继续尝试发送测试页
		fmt.Println("⚠️ 继续尝试发送测试页...")
		return sendTestPage("8443", localPort) // 使用默认端口8443
	} else {
		fmt.Printf("✅ Clodop服务响应正常 (端口: %d)\n", clodopPort)

		// 如果Clodop服务可用，尝试发送测试页
		return testClodopService(clodopPort)
	}
}

// detectClodopPort 智能检测Clodop服务端口
func detectClodopPort(userPort int) (int, error) {
	// 端口检测优先级：用户配置端口 → 8443 → 8000 → 8080 → 9000
	testPorts := []int{userPort, 8443, 8000, 8080, 9000}

	// 创建跳过SSL验证的HTTP客户端
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   3 * time.Second,
	}

	for _, port := range testPorts {
		if port <= 0 {
			continue
		}

		fmt.Printf("🔍 尝试端口 %d...\n", port)

		// 尝试HTTPS和HTTP
		protocols := []string{"https", "http"}
		urls := []string{
			"/CLodopfuncs.js?priority=1",
			"/CLodopfuncs.js",
			"/c_webskt/",
		}

		for _, protocol := range protocols {
			for _, url := range urls {
				testURL := fmt.Sprintf("%s://localhost:%d%s", protocol, port, url)

				resp, err := client.Get(testURL)
				if err == nil && resp != nil {
					resp.Body.Close()
					if resp.StatusCode == 200 {
						fmt.Printf("✅ 发现Clodop服务: %s\n", testURL)
						return port, nil
					}
				}
			}
		}
	}

	return 0, fmt.Errorf("未找到可用的Clodop服务端口，尝试了端口: %v", testPorts)
}

// testClodopService 测试Clodop服务
func testClodopService(localPort int) error {
	// 智能检测Clodop服务端口
	clodopPort, err := detectClodopPort(localPort)
	if err != nil {
		return fmt.Errorf("Clodop服务检测失败: %v", err)
	}

	fmt.Printf("✅ Clodop服务响应正常 (端口: %d)\n", clodopPort)

	// 创建测试页面进行真实打印
	err = createTestPrintPage(clodopPort)
	if err != nil {
		return fmt.Errorf("创建打印测试页失败: %v", err)
	}

	fmt.Println("📄 已创建测试打印页面，即将在浏览器中打开...")

	// 等待一秒确保文件写入完成
	time.Sleep(1 * time.Second)

	// 使用系统默认浏览器打开测试页面
	return runCommand("open", "/tmp/clodop_test.html")
}

// createTestPrintPage 创建测试打印页面
func createTestPrintPage(port int) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>HPRT打印机测试</title>
    <meta charset="UTF-8">
</head>
<body>
    <h2>HPRT打印机连接测试</h2>
    <div id="log"></div>
    <button onclick="startTest()">开始打印测试</button>
    
    <script>
        var logDiv = document.getElementById('log');
        var scriptUrls = [
            'https://localhost:%d/CLodopfuncs.js?priority=1',
            'https://localhost:%d/CLodopfuncs.js',
            'http://localhost:%d/CLodopfuncs.js?priority=1',
            'http://localhost:%d/CLodopfuncs.js'
        ];
        var currentUrlIndex = 0;
        
        function addLog(msg) {
            var time = new Date().toLocaleTimeString();
            logDiv.innerHTML += time + ': ' + msg + '<br>';
            console.log(time + ': ' + msg);
        }
        
        function loadNextScript() {
            if (currentUrlIndex >= scriptUrls.length) {
                addLog('❌ 所有Clodop脚本加载失败');
                return;
            }
            
            var url = scriptUrls[currentUrlIndex];
            addLog('🔄 尝试加载: ' + url);
            
            var script = document.createElement('script');
            script.src = url;
            script.onload = function() {
                addLog('✅ 脚本加载成功: ' + url);
                checkClodop();
            };
            script.onerror = function() {
                addLog('❌ 脚本加载失败: ' + url);
                currentUrlIndex++;
                setTimeout(loadNextScript, 500);
            };
            document.head.appendChild(script);
        }
        
        function checkClodop() {
            try {
                // 尝试不同的获取方式
                var LODOP = null;
                
                if (typeof getCLodop === 'function') {
                    LODOP = getCLodop();
                    addLog('✅ 通过getCLodop()获取到CLODOP对象');
                } else if (typeof window.CLODOP !== 'undefined') {
                    LODOP = window.CLODOP;
                    addLog('✅ 通过window.CLODOP获取到对象');
                } else if (typeof window.LODOP !== 'undefined') {
                    LODOP = window.LODOP;
                    addLog('✅ 通过window.LODOP获取到对象');
                } else {
                    addLog('❌ 未找到CLODOP对象');
                    return;
                }
                
                if (LODOP && typeof LODOP.PRINT_INIT === 'function') {
                    addLog('✅ CLODOP对象验证成功，打印功能可用');
                    addLog('📋 版本信息: ' + (LODOP.VERSION || '未知'));
                } else {
                    addLog('❌ CLODOP对象无效或缺少打印函数');
                }
                
            } catch (e) {
                addLog('❌ 检查CLODOP对象时出错: ' + e.message);
            }
        }
        
        function startTest() {
            try {
                var LODOP = null;
                
                // 获取CLODOP对象
                if (typeof getCLodop === 'function') {
                    LODOP = getCLodop();
                } else if (typeof window.CLODOP !== 'undefined') {
                    LODOP = window.CLODOP;
                } else if (typeof window.LODOP !== 'undefined') {
                    LODOP = window.LODOP;
                }
                
                if (!LODOP) {
                    addLog('❌ 无法获取CLODOP对象');
                    return;
                }
                
                addLog('🖨️ 开始执行打印测试...');
                
                // 初始化打印任务
                LODOP.PRINT_INIT("HPRT测试页");
                
                // 设置纸张大小 (热敏打印机通常用80mm宽度)
                LODOP.SET_PRINT_PAGESIZE(1, 0, 0, "80mm*120mm");
                
                // 添加打印内容
                LODOP.ADD_PRINT_TEXT(50, 10, 200, 30, "HPRT打印机测试");
                LODOP.ADD_PRINT_TEXT(100, 10, 300, 20, "时间: %s");
                LODOP.ADD_PRINT_TEXT(130, 10, 300, 20, "状态: 打印机工作正常");
                LODOP.ADD_PRINT_TEXT(160, 10, 300, 20, "配置: 网络连接已建立");
                LODOP.ADD_PRINT_TEXT(190, 10, 300, 20, "测试: 端口%d通信正常");
                
                // 执行打印
                var result = LODOP.PRINT();
                
                if (result) {
                    addLog('✅ 打印任务已发送，任务ID: ' + result);
                    addLog('🎉 测试完成！请查看打印机输出');
                } else {
                    addLog('❌ 打印任务发送失败');
                }
                
            } catch (e) {
                addLog('❌ 打印测试失败: ' + e.message);
            }
        }
        
        // 页面加载完成后开始加载脚本
        window.onload = function() {
            addLog('📄 页面加载完成，开始加载Clodop脚本');
            loadNextScript();
        };
    </script>
</body>
</html>`, port, port, port, port, now, port)

	return os.WriteFile("/tmp/clodop_test.html", []byte(html), 0644)
}

// sendTestPage 发送测试打印页
func sendTestPage(clodopPort, localPort string) error {
	// 通过JavaScript命令调用Clodop
	// 动态检测协议和端口
	testHTML := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>HPRT打印测试</title>
</head>
<body>
    <div style="padding: 20px; font-family: Arial, sans-serif;">
        <h2>HPRT打印机测试页面</h2>
        <p>正在测试Clodop打印功能...</p>
        <div id="status">正在加载...</div>
        <div id="debug" style="margin-top: 20px; padding: 10px; background-color: #f0f0f0; font-family: monospace; font-size: 12px;"></div>
        <br>
        <p><strong>如果看到这个页面，说明:</strong></p>
        <ul>
            <li>✓ 网络连接正常</li>
            <li>✓ 端口转发工作正常</li>
        </ul>
        <p><em>注意: 请确保在远程Windows电脑上已安装并运行Clodop服务</em></p>
    </div>

    <script type="text/javascript">
        // 调试信息
        function addDebug(msg) {
            var debug = document.getElementById('debug');
            debug.innerHTML += new Date().toLocaleTimeString() + ': ' + msg + '<br>';
        }

        // 动态加载Clodop脚本
        function loadClodopScript() {
            var port = '` + clodopPort + `';
            var urls = [
                'https://localhost:' + port + '/CLodopfuncs.js?priority=1',
                'https://localhost:' + port + '/CLodopfuncs.js',
                'http://localhost:' + port + '/CLodopfuncs.js',
                'http://localhost:' + port + '/CLodopfuncs'
            ];

            var tryIndex = 0;

            function tryNextUrl() {
                if (tryIndex >= urls.length) {
                    document.getElementById('status').innerHTML = '<span style="color: red;">❌ 无法加载Clodop脚本，所有URL都失败了</span>';
                    addDebug('所有Clodop URL都加载失败');
                    return;
                }

                var url = urls[tryIndex];
                addDebug('尝试加载: ' + url);
                
                var script = document.createElement('script');
                script.type = 'text/javascript';
                script.src = url;
                
                script.onload = function() {
                    addDebug('脚本加载成功: ' + url);
                    setTimeout(testPrint, 1000); // 等待1秒再测试打印
                };
                
                script.onerror = function() {
                    addDebug('脚本加载失败: ' + url);
                    tryIndex++;
                    setTimeout(tryNextUrl, 500); // 等待0.5秒再尝试下一个
                };
                
                document.head.appendChild(script);
            }

            tryNextUrl();
        }

        // 测试打印功能
        function testPrint() {
            try {
                addDebug('检查getLodop函数...');
                if (typeof getLodop === 'undefined') {
                    document.getElementById('status').innerHTML = '<span style="color: red;">❌ getLodop函数未定义</span>';
                    addDebug('getLodop函数未定义');
                    return;
                }

                addDebug('调用getLodop()...');
                var LODOP = getLodop();
                if (LODOP) {
                    addDebug('LODOP对象获取成功');
                    
                    LODOP.PRINT_INIT("HPRT测试页");
                    LODOP.SET_PRINT_PAGESIZE(1, 0, 0, "80mm*120mm");
                    
                    // 添加标题
                    LODOP.ADD_PRINT_TEXT(20, 50, 200, 30, "HPRT打印机测试页");
                    LODOP.SET_PRINT_STYLEA(0, "FontName", "微软雅黑");
                    LODOP.SET_PRINT_STYLEA(0, "FontSize", 14);
                    LODOP.SET_PRINT_STYLEA(0, "Bold", 1);
                    
                    // 添加分隔线
                    LODOP.ADD_PRINT_TEXT(50, 20, 240, 20, "================================");
                    
                    // 添加状态信息
                    LODOP.ADD_PRINT_TEXT(80, 50, 200, 20, "✓ 打印机工作正常！");
                    LODOP.SET_PRINT_STYLEA(0, "FontSize", 12);
                    
                    LODOP.ADD_PRINT_TEXT(110, 50, 200, 20, "✓ VPN连接正常");
                    LODOP.ADD_PRINT_TEXT(130, 50, 200, 20, "✓ 端口转发正常");
                    LODOP.ADD_PRINT_TEXT(150, 50, 200, 20, "✓ 网络通信正常");
                    
                    // 添加时间信息
                    var now = new Date();
                    var timeStr = now.getFullYear() + "-" + 
                                 (now.getMonth()+1).toString().padStart(2,'0') + "-" + 
                                 now.getDate().toString().padStart(2,'0') + " " +
                                 now.getHours().toString().padStart(2,'0') + ":" + 
                                 now.getMinutes().toString().padStart(2,'0') + ":" + 
                                 now.getSeconds().toString().padStart(2,'0');
                    
                    LODOP.ADD_PRINT_TEXT(180, 50, 200, 20, "测试时间: " + timeStr);
                    LODOP.SET_PRINT_STYLEA(0, "FontSize", 10);
                    
                    // 添加配置信息
                    LODOP.ADD_PRINT_TEXT(210, 50, 200, 20, "配置工具版本: 1.0");
                    
                    // 添加结束分隔线
                    LODOP.ADD_PRINT_TEXT(240, 20, 240, 20, "================================");
                    
                    addDebug('准备执行打印...');
                    // 执行打印
                    LODOP.PRINT();
                    
                    // 显示成功信息
                    document.getElementById('status').innerHTML = '<span style="color: green;">✅ 打印命令已发送，请检查打印机是否出纸</span>';
                    addDebug('打印命令执行成功');
                } else {
                    document.getElementById('status').innerHTML = '<span style="color: red;">❌ 无法获取LODOP对象</span>';
                    addDebug('无法获取LODOP对象');
                }
            } catch (e) {
                document.getElementById('status').innerHTML = '<span style="color: red;">❌ 打印测试失败: ' + e.message + '</span>';
                addDebug('异常: ' + e.message);
            }
        }

        // 页面加载完成后开始测试
        window.onload = function() {
            addDebug('页面加载完成，开始加载Clodop脚本');
            loadClodopScript();
        };
    </script>
</body>
</html>`

	// 创建临时HTTP服务器来提供测试页面
	server := &http.Server{
		Addr: ":0", // 使用随机端口
	}

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(testHTML))
	})

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("无法启动测试服务器: %v", err)
	}
	defer listener.Close()

	testPort := listener.Addr().(*net.TCPAddr).Port
	testURL := fmt.Sprintf("http://localhost:%d/test", testPort)

	// 启动服务器
	go server.Serve(listener)

	// 让系统默认浏览器打开测试页面
	fmt.Printf("📱 打开浏览器测试页面: %s\n", testURL)

	// 在macOS上打开浏览器
	err = runCommand("open", testURL)
	if err != nil {
		return fmt.Errorf("无法打开浏览器: %v", err)
	}

	// 等待更长时间让页面加载和执行
	time.Sleep(8 * time.Second)

	return nil
}

// testLocalPort 测试本地端口是否可用
func testLocalPort(port string) error {
	conn, err := net.DialTimeout("tcp", "localhost:"+port, 5*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接到本地端口 %s: %v", port, err)
	}
	defer conn.Close()
	return nil
}

// testRemoteConnection 测试远程连接
func testRemoteConnection(host, port string) error {
	conn, err := net.DialTimeout("tcp", host+":"+port, 10*time.Second)
	if err != nil {
		return fmt.Errorf("无法连接到远程主机 %s:%s: %v", host, port, err)
	}
	defer conn.Close()
	return nil
}
