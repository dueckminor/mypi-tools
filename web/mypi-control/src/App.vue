<template>
  <iframe ref="iframe" width="100%" height="100%"></iframe>
</template>

<style>
body, html
{
    margin: 0; padding: 0; height: 100%; overflow: hidden;
}
#content
{
    position:absolute; left: 0; right: 0; bottom: 0; top: 0px; 
}
</style>


<script>
export default {
  name: 'app',
  data: function () {
    return {
    }
  },
  methods: {
  },
  mounted: function() {
    var rnd = Math.floor(Math.random()*100000) + 1
    this.$refs.iframe.src = "https://nodered.rh94.dueckminor.de/nodered/ui?time="+rnd;
  }
}
</script>
