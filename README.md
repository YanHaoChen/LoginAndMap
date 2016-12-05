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

## 資料儲存

### 資料庫
* name：go_db
* table:
	* LightMapMarker
	* User

#### User

NAME | Type
--- | ---
id | bigint(20)
account | varchar(255)
password | varchar(255)
token | varchar(255)

#### LightMapMarker

NAME | Type
--- | ---
id | bigint(20)
Time | datetime
value | double
account | varchar(255)
locationName | varchar(255)

### SessionStorage

NAME | 用途
--- | ---
account | 驗證
token | 驗證
find | 給 highchart.html 使用，儲存顯示地理位置

## 流程

1. 登入
2. 地圖及 Marker 呈現
3. 單擊顯示地理名稱
4. 雙擊顯示此地的數據圖表呈現

### 登入

取出帳號、密碼，使用 API （/User/Login）進行驗證。成功則將`token`及`account`存入 SessionStorage 供後續驗證使用。

##### 登入（login.html）

```html
$("#login").click(function() {
    $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8080/User/Login",
        dataType: "json",
        data: {
            account: $("#account").val(),
            password: $("#password").val()
        },
        success: function(data) {
            sessionStorage.setItem("token",data["token"]);
            sessionStorage.setItem("account",data["account"]);

            location.href="/hw3_web/marksmap.html";
        },
        error: function(jqXHR) {
            alert("發生錯誤: " + jqXHR.status);
        }
    })
  })
```

##### 登入（server.go）

```go
func Login(c *gin.Context) {
  var user User
	account := c.PostForm("account")
  password := c.PostForm("password")
	if account != "" && password != "" {

		err := dbmap.SelectOne(&user,"SELECT * FROM User WHERE account=? AND password=?", account, password)
			if err == nil {
				b := make([]byte, 10)
				rand.Read(b)
				str := fmt.Sprintf("%x", b)
				user.Token = str
				dbmap.Update(&user)

				c.JSON(200, gin.H{"token": str,"account": user.Account})
			} else {

				c.JSON(404, gin.H{"error": err})
			}
		} else {
		c.JSON(422, gin.H{"error": "fields are empty"})
	}
}
```


### 地圖及 Marker 呈現

驗證（User/AuthStatus）後，存取對應 Marker（Marker/GetLightMapMarkers），並在 Marker 上加入`click`及`dblclick`兩項監聽事件。

##### 驗證（marksmap.html）

```html
$(document).ready(function(){
      $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8080/User/AuthStatus",
        dataType: "json",
        data: {
          account: sessionStorage["account"],
          token: sessionStorage["token"]
        },
        success: function(data){
          getTheMarkers();
        },
        error: function(jqXHR) {
          sessionStorage["account"] = "";
          sessionStorage["token"] = "";
          location.href="/hw3_web/login.html";

        }
      });
    });
```

##### 驗證（server.go）

```go
func AuthStatus(c *gin.Context) {
	var user User
	account := c.PostForm("account")
	token := c.PostForm("token")

	if account != "" && token != "" {
		err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE account=? AND token=?", account, token)
		if err == nil {
			c.JSON(200,gin.H{"status":"pass"})
		} else {
			c.JSON(404, gin.H{"error": err})
		}
	} else {
		c.JSON(422, gin.H{"error":"fields are empty"})
	}
}
```

##### 存取對應 Marker（marksmap.html）

```go
function getTheMarkers(){
      $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8080/Marker/GetLightMapMarkers",
        dataType: "json",
        data: {
          account: sessionStorage["account"],
          token: sessionStorage["token"]
        },
        success: function(data){
          makers = data;
          addMarkers();
        },
        error: function(jqXHR) {
          sessionStorage["account"] = "";
          sessionStorage["token"] = "";
          location.href="/hw3_web/login.html";
        }
      });
    }
```
##### 存取對應 Marker（server.go）

```go
func GetLightMapMarkers(c *gin.Context) {
	var user User
	var markers []LightMapMarker
	account := c.PostForm("account")
	token := c.PostForm("token")

	if token != "" && account != "" {
		err := dbmap.SelectOne(&user, "SELECT * FROM User WHERE token=? AND account=?", token,account)
		if err == nil {
			_, err := dbmap.Select(&markers, "SELECT * FROM LightMapMarker WHERE account=?",account)
			if err == nil {
				c.JSON(200,markers)
			} else {
				c.JSON(404,gin.H{"error":err.Error()})
			}

		} else {
			c.JSON(404,gin.H{"error":"user is wrong"})
		}
	} else {
		c.JSON(422, gin.H{"error":"fields are empty"})
	}

}
```

##### 加入 Marker 並加入監聽事件（marksmap.html）

```html
function addMarkers(){
      var ads = [];
      for(var i=0;i< makers.length; i++){
        var locationName = makers[i]["locationName"];

        if($.inArray(locationName.toString(),ads) == -1) {
          ads.push(locationName.toString());
          var geocoder = new google.maps.Geocoder();
          geocoder.geocode({address: locationName}, function(results, status){
            if (status == google.maps.GeocoderStatus.OK) {
              var marker = new google.maps.Marker({
                position:results[0].geometry.location,
                map: map,
                animation: google.maps.Animation.DROP,
                title: locationName
              });
              marker.addListener('click', function(e) {
                html = "";
                html += locationName;
                var infoWindow = new google.maps.InfoWindow({
                  content: html
                });
                infoWindow.open(map,marker);
              });
              marker.addListener('dblclick', function(e) {
                sessionStorage["find"] = locationName.toString();
                window.open('highchart.html', locationName.toString(), config='height=400,width=600')
              });
            }
          });

```
### 單擊顯示地理名稱

在單擊監聽事件中，加入開啟 InfoWindow。

##### 加入開啟 InfoWindow（marksmap.html）

```html
marker.addListener('click', function(e) {
                html = "";
                html += locationName;
                var infoWindow = new google.maps.InfoWindow({
                  content: html
                });
                infoWindow.open(map,marker);
              });
```

### 雙擊顯示此地的數據圖表呈現

在雙擊監聽事件中，開啟另一個 html（highchart.html）。highchart.html 再透過 API 取得特定地點資料，並透過 highchart 呈現。

##### 在雙擊監聽事件中，開啟另一個 html（marksmap.html）

```html
marker.addListener('dblclick', function(e) {
                sessionStorage["find"] = locationName.toString();
                window.open('highchart.html', locationName.toString(), config='height=400,width=600')
              });
```
##### 取得特定地點資料（highchart.html）
```html
function initialize(){
      $.ajax({
        type:'POST',
        url:'http://127.0.0.1:8080/Marker/GetMarkerValue',
        data:{
          locationName: sessionStorage["find"],
          account: sessionStorage["account"]
        },
        dataType: 'json',
        success:getDataSuccess
      })
    };
```

##### 取得特定地點資料（server.go）

```go
func GetMarkerValue(c *gin.Context) {
	var markers []LightMapMarker
	account := c.PostForm("account")
	locationName := c.PostForm("locationName")

	_, err := dbmap.Select(&markers, "SELECT * FROM LightMapMarker WHERE account=? AND locationName=?",account, locationName)
	if err == nil {
		c.JSON(200,markers)
	} else {
		c.JSON(404,gin.H{"error":"something wrong"})
	}
}
```

## 完成畫面

### 登入
![登入](https://github.com/YanHaoChen/LoginAndMap/blob/master/images/login.png?raw=true)

### 地圖
![地圖](https://github.com/YanHaoChen/LoginAndMap/blob/master/images/map.png?raw=true)

### 顯示數值
![顯示數值](https://github.com/YanHaoChen/LoginAndMap/blob/master/images/map.png?raw=true)

### Server 狀況
![Server 狀況](https://github.com/YanHaoChen/LoginAndMap/blob/master/images/server.png?raw=true)