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
        $files = scandir("./files", SCANDIR_SORT_ASCENDING);
    ?>
    <img id="background" src="/static/img/background_1.png">
    <div class="center" style="overflow-y: scroll">
        <?php
            foreach ($files as $file) {
                if ($file === "." || $file === "..") {
                    continue;
                }
                $base = basename($file, ".html");
                echo "<p><a href='files/$base'>$base</a></p>\n";
            }
        ?>
    </div>
</body>
</html>