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
</style>
<body>
    <img id="background" src="background_1.png">
    <script src="https://cdn.jsdelivr.net/npm/vue@2.5.13/dist/vue.js"></script>
    Today's RFC Number is: <h1>
    <?php
        $content = file_get_contents(".rfc");
        list($num, $desc) = explode(":", $content);
        print($num);
    ?>
    </h1>
    With description:
    <h1>
        <?php
            print($desc);
        ?>
    </h1>
</body>
</html>
