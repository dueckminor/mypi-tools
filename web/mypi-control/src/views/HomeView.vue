<template>
  <v-container>
    <v-card>
      <v-card-title>
        <v-btn color="secondary" fab x-small dark @click="powerClicked">
          <v-icon>mdi-power</v-icon> </v-btn
        >&nbsp;&nbsp;RPI
      </v-card-title>
      <v-row>
        <v-col> </v-col>
        <v-col>
          <GChart
            :settings="{ packages: ['corechart', 'gauge'] }"
            type="Gauge"
            :data="gaugeData"
            :options="gaugeOptions"
          />
        </v-col>
      </v-row>
    </v-card>
  </v-container>
</template>

<script>
// @ is an alias to /src
//import Vue from 'vue'
import { GChart } from "vue-google-charts";

export default {
  name: "HomeView",
  data: () => ({
    gaugeData: [
      ["",""],
      ["CPU", 0],
      ["RAM", 0],
      //["SWAP", 0],
    ],
    gaugeOptions: {
      width: 400,
      height: 120,
      redFrom: 90,
      redTo: 100,
      yellowFrom: 75,
      yellowTo: 90,
      minorTicks: 5,
    },
  }),
  components: {
    GChart,
  },
  mounted() {
    this.sockets.subscribe("stats/rpi/cpu", (data) => {
      this.gaugeData[1]=["CPU",parseFloat(data)];
    });
    this.sockets.subscribe("stats/rpi/mem", (data) => {
      this.gaugeData[2]=["RAM",parseFloat(data)];
    });
    this.sockets.subscribe("stats/rpi/swap", (data) => {
      this.gaugeData[3]=["SWAP",parseFloat(data)];
    });
  },
  methods: {
    powerClicked: function() {
      this.$socket.emit("notice", "foo");
    },
  },
};
</script>
