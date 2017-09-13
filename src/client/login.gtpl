<!DOCTYPE html>
<html>
    <head>
        <title></title>
        
        <meta charset="utf8"/>
        
        <style>
            #container{
                position:relative;
                
                margin-left:10%;
                width:80%;
                
                text-align: center;
            }

            #sell_form{
                position: relative;
                
                display:inline-block;
                
                padding:20px;

                border:1px solid #777;
                border-radius:5px;
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

            input[type=submit]{
                background: #CCC;
                padding:5px 0 5px 0;
            }

            td{
                padding:5px;
            }

            td:nth-child(1){
                text-align: right;
            }
        </style>

        <script>
        </script>
    </head>

    <body>
        <div id="container">
            <form id="sell_form" action="/sell_form" method="post">
                <table>
                    <tr>
                        <td>
                            Login:
                        </td>

                        <td>
                            <input type="text" name="login"/>                    
                        </td>
                    </tr>

                    <tr>
                        <td>
                            Password:        
                        </td>

                        <td>
                            <input type="password" name="password"/>
                        </td>
                    </tr>
                    
                    <tr>
                        <td>
                        </td>
                        
                        <td>
                            <input type="submit" value="Login"/>
                        </td>
                    </tr>             
                </table>
            </form>
        </div>
    </body>
</html>