<!DOCTYPE html>
<html>
    <head>
        <title></title>
        
        <meta charset="utf8"/>
        
        <style>
            #container{
                position:relative;

                width:80%;
                
                margin-left:10%;
                
                text-align: center;
            }
            
            #history{
                position:relative;
                
                display:inline-block;

                text-align: center;
            }

            td{
                padding:0 10px 0 10px;

                border:1px solid #777;
            }
        </style>
    </head>

    <body>
        <div id="container">
            <div id="history">
                <a href="/sell_form">Назад</a>

                <p>

                <table>
                    <tr>
                        <td>ID</td>
                        <td>Time</td>
                        <td>UserID</td>
                        <td>Ticker</td>
                        <td>Amount</td>
                        <td>Price</td>
                    </tr>
                    {{range $k, $v := .}}
                    {{if eq .Status ""}}
                    <tr>
                        <td>{{$v.ID}}</td>
                        <td>{{$v.Time}}</td>
                        <td>{{$v.UserID}}</td>
                        <td>{{$v.Ticker}}</td>
                        <td>{{$v.Amount}}</td>
                        <td>{{$v.Price}}</td>
                    </tr>
                    {{end}}
                    {{end}}
                </table>
            </div>
            
            <p>

            <a href="/sell_form">Назад</a>
        </div>
    </body>
</html>

<!-- type Deal struct {
    ID     int     `json:"-"`
    Time   int     `json:"time"`
    UserID int     `json:"-"`
    Ticker string  `json:"ticker"`
    Amount int     `json:"amount"`
    Price  float64 `json:"price"`
    Status string  `json:"status"`
} -->