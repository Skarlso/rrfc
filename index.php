<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Random RFC</title>
</head>
<style>
    img#background {
      width: 100%;
      height: auto;
      position: fixed;
      z-index: -1;
      bottom: 0;
      right: 0;
    }
    .center {
        height: auto;
        margin: auto;
        padding: 10px;
        text-align: center;
    }
</style>
<body>
    <img id="background" src="background_1.png">
    <div class="center">
        <h1>
            <?php
                $content = file_get_contents(".rfc");
                list($num, $desc) = explode(":", $content);
                print($num);
            ?>
        </h1>
        <h1>
            <?php
                print($desc);
            ?>
        </h1>
    </div>
    <div id="disqus_thread"></div>
    <script>
        /**
         *  RECOMMENDED CONFIGURATION VARIABLES: EDIT AND UNCOMMENT
         *  THE SECTION BELOW TO INSERT DYNAMIC VALUES FROM YOUR
         *  PLATFORM OR CMS.
         *
         *  LEARN WHY DEFINING THESE VARIABLES IS IMPORTANT:
         *  https://disqus.com/admin/universalcode/#configuration-variables
         */
        var disqus_config = function () {
            // Replace PAGE_URL with your page's canonical URL variable
            this.page.url = "https://rrfc.app";

            // Replace PAGE_IDENTIFIER with your page's unique identifier variable
            this.page.identifier = "1234asdf";
        };

        (function() {  // REQUIRED CONFIGURATION VARIABLE: EDIT THE SHORTNAME BELOW
            var d = document, s = d.createElement('script');

            // IMPORTANT: Replace EXAMPLE with your forum shortname!
            s.src = 'https://hannibalDisqus.disqus.com/embed.js';

            s.setAttribute('data-timestamp', +new Date());
            (d.head || d.body).appendChild(s);
        })();
    </script>
    <noscript>
        Please enable JavaScript to view the
        <a href="https://disqus.com/?ref_noscript" rel="nofollow">
            comments powered by Disqus.
        </a>
    </noscript>
</body>
</html>
