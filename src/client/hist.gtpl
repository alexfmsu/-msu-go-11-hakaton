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
                        <td>Interval</td>
                        <td>Open</td>
                        <td>High</td>
                        <td>Low</td>
                        <td>Close</td>
                        <td>Ticker</td>
                    </tr>
                    {{range $k, $v := .}}
                    <tr>
                        <td>{{$v.ID}}</td>
                        <td>{{$v.Time}}</td>
                        <td>{{$v.Interval}}</td>
                        <td>{{$v.Open}}</td>
                        <td>{{$v.High}}</td>
                        <td>{{$v.Low}}</td>
                        <td>{{$v.Close}}</td>
                        <td>{{$v.Ticker}}</td>
                    </tr>
                    {{end}}
                </table>
            </div>
            
            <p>

            <a href="/sell_form">Назад</a>
        </div>
    </body>
</html>

<!--
type OHLCV struct {
    ID       int64   `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
    Time     int32   `protobuf:"varint,2,opt,name=Time" json:"Time,omitempty"`
    Interval int32   `protobuf:"varint,3,opt,name=Interval" json:"Interval,omitempty"`
    Open     float32 `protobuf:"fixed32,4,opt,name=Open" json:"Open,omitempty"`
    High     float32 `protobuf:"fixed32,5,opt,name=High" json:"High,omitempty"`
    Low      float32 `protobuf:"fixed32,6,opt,name=Low" json:"Low,omitempty"`
    Close    float32 `protobuf:"fixed32,7,opt,name=Close" json:"Close,omitempty"`
    Ticker   string  `protobuf:"bytes,8,opt,name=Ticker" json:"Ticker,omitempty"`
}
-->
