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
        </style>
    </head>

    <body>
        <div id="container">
            {{.error}}

            <p>

            <a href={{.back}}>Назад</a>
        </div>
    </body>
</html>
