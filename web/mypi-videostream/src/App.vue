<template>
  <v-app light id="app">
    <v-toolbar app>
      <v-btn @click.stop="folder()" >
        <v-icon>folder_open</v-icon>
      </v-btn>
      <v-spacer />
      <v-btn @click.stop="prev()" class="hidden-sm-and-down">
        <v-icon>skip_previous</v-icon>
      </v-btn>
      <v-btn @click.stop="play()">
        <v-icon>play_arrow</v-icon>
      </v-btn>
      <v-btn @click.stop="pause()" class="hidden-sm-and-down">
        <v-icon>pause</v-icon>
      </v-btn>
      <v-btn @click.stop="next()" class="hidden-sm-and-down">
        <v-icon>skip_next</v-icon>
      </v-btn>
      <v-spacer />
      <v-btn @click.stop="fullscreen()" >
        <v-icon>fullscreen</v-icon>
      </v-btn>
    </v-toolbar>
    <v-content>
      <video ref="video" muted="muted" autoplay playsinline width="100%"/>
    </v-content>
  </v-app>
</template>

<script>
import VueRouter from 'vue-router'

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
