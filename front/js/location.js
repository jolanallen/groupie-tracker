document.addEventListener("DOMContentLoaded", function () {
  const iframe = document.getElementById("map-iframe");
  iframe.onload = function () {
    iframe.contentWindow.scrollTo(0, iframe.contentWindow.document.body.scrollHeight);
  };
});
