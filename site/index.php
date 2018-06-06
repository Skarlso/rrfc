<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Random RFC</title>
</head>
<link rel="stylesheet" type="text/css" href="/static/css/main.css">
<script src="/static/js/disqus.js"></script>
<body>
    <?php
        $content = file_get_contents(".rfc");
        list($num, $desc) = explode(":", $content);
    ?>
    <img id="background" src="/static/img/background_1.png">
    <div class="center">
        <h1 id="number">
            <?php print($num); ?>
        </h1>
        <h1 id="description">
            <?php print($desc); ?>
        </h1>
    </div>
    <div class="center" id="disqus_thread"></div>
    <noscript>
        Please enable JavaScript to view the
        <a href="https://disqus.com/?ref_noscript" rel="nofollow">
            comments powered by Disqus.
        </a>
    </noscript>
</body>
</html>