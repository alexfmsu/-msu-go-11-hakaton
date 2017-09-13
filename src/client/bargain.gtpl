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
                border-radius: 5px;
            }

            input:focus {
                border: 1px solid #39c;
            }

            #deal_form{
                display:inline-block;
                width:25%;

                text-align: center;
            }

            #deal_form input{
                width:100%;
            
                padding-left: 3px;
            }

            #deal_form select{
                width:100%;
            }

            nav{
                padding-bottom:40px;
            }
            nav a{
                padding:0 10 0 10px;
            }

            td{
                padding:5px 5px 5px 5px;
            }
        </style>

    </head>
    
    <body>
        <div id="container">
            <a href = "/sell_form">Продать</a>
            <a href = "/buy_form">Купить</a>
            
            <nav>
                <a href="/me">История сделки</a>
                <a href="/deals">Новая сделка</a>
                <a href="/bargain">История сделки</a>
            </nav>
            
            <p>
            
            <div id="sell_form">

                <form action="/sell" method="get" id="deal_form">
                    <table>
                        <tr>
                            <td>
                                <input type="text" name="amount" placeholder="количество">
                            </td>
                        </tr>

                        <tr>
                            <td>
                                <input type="price" name="price" placeholder="цена">
                            </td>
                        </tr>
                        
                        <tr>
                            <td>
                                <select name="ticker">
                                    <option>
                                        Rim7
                                    </option>
                                    <option>
                                        Инструмент1
                                    </option>
                                    <option>
                                        Инструмент2
                                    </option>
                                </select>
                            </td>
                        </tr>
                        
                        <tr>
                            <td>
                                <input type="submit" value="Продать" class="btn"/>
                            </td>
                        </tr>
                    </table>
                </form>
            </div>
        </div>
    </body>
</html> 
