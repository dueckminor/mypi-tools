<template>
  <v-app light id="app">
    <v-app-bar app>
      <v-btn @click.stop="folder()" >
        <v-icon >mdi-folder-open</v-icon>
      </v-btn>
      <v-spacer />
      <v-btn @click.stop="prev()" class="hidden-sm-and-down">
        <v-icon>mdi-skip-previous</v-icon>
      </v-btn>
      <v-btn @click.stop="play()">
        <v-icon>mdi-play</v-icon>
      </v-btn>
      <v-btn @click.stop="pause()" class="hidden-sm-and-down">
        <v-icon>mdi-pause</v-icon>
      </v-btn>
      <v-btn @click.stop="next()" class="hidden-sm-and-down">
        <v-icon>mdi-skip-next</v-icon>
      </v-btn>
      <v-spacer />
      <v-btn @click.stop="fullscreen()" >
        <v-icon>mdi-fullscreen</v-icon>
      </v-btn>
    </v-app-bar>
    <v-main>
      <video ref="video" muted="muted" autoplay playsinline :style="{'width': '100%', 'height':'100%'}"/>
    </v-main>
  </v-app>
</template>

<script>
import Hls from "hls.js";

export default {
  name: 'app',
  components: {
  },
  data: function () {
    return {
    }
  },
  computed: {
    videoElement () {
      return this.$refs.video;
    }, 
  },
  methods: {
    play() {
      var video = this.videoElement
      if (video.paused) { 
        video.play(); 
      } else if (Hls.isSupported()) {
          var hls = new Hls();
          hls.loadSource("/cams/0/stream.m3u8");
          hls.attachMedia(this.videoElement);
          hls.on(Hls.Events.MANIFEST_PARSED, function () {
              this.videoElement.play();
          })
      } else {
        var source = document.createElement('source');
        source.setAttribute('src', '/cams/0/stream.m3u8')
        video.innerHTML = '';
        video.appendChild(source)
        video.play();
      }
    },
    pause() {
      var video = this.videoElement
      video.pause()
    },
    fullscreen() {
      var video = this.videoElement
      if (video.requestFullscreen) {
          video.requestFullscreen();
      } else if (video.mozRequestFullScreen) {
          video.mozRequestFullScreen();
      } else if (video.msRequestFullscreen) {
          video.msRequestFullscreen();
      } else if (video.webkitRequestFullscreen) {
          video.webkitRequestFullscreen();
      } else if (video.webkitEnterFullscreen) {
        video.webkitEnterFullscreen();
      }
    }
  },
  beforeMount() {
  },
  mounted() {
    this.play()
  }
}
</script>

<style>
</style>
