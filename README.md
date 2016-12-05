# LoginAndMap

## 整體架構
* 後端：gin Framework（golang）
* 前端：jquery、bootstrap、Google Map API v3

## 前端功能
* 登入（login.html）：

功能 | _
--- | ---
後端驗證 | OK
密碼經雜湊傳入後端 | 
使用 token 進行後續驗證 | OK

* 地圖（marksmap.html）：

功能 | _
--- | ---
使用 token 進行後續驗證 | OK
以使用者帳號及 token 進行 marker 搜尋 | OK
marker 單擊：marker 地理名稱 | OK
marker 雙擊：此地的數據圖表呈現 | OK 
登出按鈕 | OK

* 圖表（highchart.html）:

功能 | _
--- | ---
使用 highchart 進行圖表呈現 | OK
以使用者帳號及地理位址進行 marker 搜尋 | OK

## 後端功能（Router）

### User Group

Method | 用途 | 參數 | Url 
--- | --- | --- |---
POST | 登入 |account:text、password:text | User/Login
POST | 驗證 | account:text、token:text | User/AuthStatus

### Marker Group

Method | 用途 | 參數 | Url 
--- | --- | --- |---
POST | 取得對應 Marker |account:text、token:text | Marker/GetLightMapMarkers
POST | 取得特定地點 Marker | account:text、locationName:text | Marker/GetMarkerValue

## 資料庫格式
