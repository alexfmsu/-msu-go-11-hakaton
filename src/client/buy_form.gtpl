<!DOCTYPE html>
<html>
    <head>
        <title></title>
        
        <meta charset="utf8"/>
        
        <style>
            html{
                font-family: sans-serif;
            }
            #container{
                position:relative;
                
                margin-left:10%;
                width:80%;
                
                text-align: center;
            }

            input{
                width:100%;
                height:100%;

                text-indent: 5px;
                padding:3px 0 3px 0;
                
                border:1px solid #777;
                border-radius: 3px;
            }

            input:focus {
                border: 1px solid #39c;
            }

            #deal_form{
                display:inline-block;

                text-align: center;
            }

            #deal_form input{
                width:100%;
            
                padding-left: 3px;
            }

            #deal_form select{
                width:100%;
            }

            td{
                padding:5px 5px 5px 5px;
            }

            .buy_btn, .sell_btn{
                display: inline-block;
                
                color:black;
                background:#DDD;

                text-align: center;
                text-decoration: none;
                
                padding:3px 10px 3px 10px;

                border:1px solid #777;
                border-radius: 3px;
            }

            a{
                color:black;
            }
        </style>

    </head>
    
    <body>
        <div id="container">
            <a href="/me">Мои позиции</a>
            <a href="/deals">История моих сделок</a>
            <a href="/bargain">Последняя история торгов</a>
            
            <p>

            <a href = "/buy_form" class="buy_btn">Купить</a>
            <a href = "/sell_form" class="sell_btn">Продать</a>
            
            <p>
            
            <div id="buy_form">
                <form action="/buy" method="get" id="deal_form">
                    <table>
                        <tr>
                            <td>
                                <input type="number" name="amount" placeholder="количество">
                            </td>
                        </tr>

                        <tr>
                            <td>
                                <input type="number" name="price" placeholder="цена">
                            </td>
                        </tr>
                        
                        <tr>
                            <td>
                                <select name="ticker">
                                    <option>
                                        RIM7
                                    </option>
                                    <option>
                                        SIM7
                                    </option>
                                </select>
                            </td>
                        </tr>
                        
                        <tr>
                            <td>
                                <input type="submit" value="Купить"/>
                            </td>
                        </tr>
                    </table>
                </form>
            </div>
        </div>
    </body>
</html> 
