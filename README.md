# othello-board
othello(reversi) arena for humans and AIs

# 遊戲截圖
![](https://user-images.githubusercontent.com/65079334/120734562-a9e3ce80-c51b-11eb-8590-6033cb762637.png)
![](https://user-images.githubusercontent.com/65079334/120681698-39629080-c4ce-11eb-8064-37293704bf51.png)

# 規則
規則與一般黑白棋相同，版面大小為6x6，若有一方無處可下會自動PASS，換另一方下  
雙方皆無處可下時遊戲結束，依棋子數目決定輸贏或平手  

# 使用外部AI說明
外部的AI程式須接收input，並輸出結果  
如：輸入```++++++++++++++OX++++XO++++++++++++++ 1```，輸出```Bc```  
(X表示黑方，O表示白方；1表示為黑方，2為白方)  
程式接收回傳值時，若回傳值不合法，GUI會顯示外部AI出錯，並將錯誤內容輸出至log file內  

# 下載
https://github.com/lemon37564/othello-board/releases

# 自行編譯
require go 1.16+  
先```go get fyne.io/fyne/v2```
### windows
```go build -ldflags="-H windowsgui"```  
### linux or macOS
```go build```  
