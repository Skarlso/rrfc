var disqus_config = function () {
    // Replace PAGE_URL with your page's canonical URL variable
    id = document.getElementById("number").innerText
    this.page.url = "https://rrfc.app/previous/" + id;

    // Replace PAGE_IDENTIFIER with your page's unique identifier variable
    this.page.identifier = id;
};
(function() {
    var d = document, s = d.createElement('script');
    s.src = 'https://rrfc.disqus.com/embed.js';
    s.setAttribute('data-timestamp', +new Date());
    (d.head || d.body).appendChild(s);
})();