<template>
  <v-container>
    <v-card>
      <v-card-title>
          <v-btn color="secondary" fab x-small dark>
            <v-icon>mdi-power</v-icon>
          </v-btn>&nbsp;&nbsp;RPI
      </v-card-title>
      <v-row>
        <v-col>

        </v-col>
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
import { GChart } from "vue-google-charts";

export default {
  name: "Home",
  data: () => ({
    gaugeData: [
      ["Label", "Value"],
      ["CPU", 55],
      ["RAM", 95],
      ["SWAP", 5],
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
};
</script>
