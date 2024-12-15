<template>
  <div>
    <h1 align="center">All Timetables</h1>
    <div v-if="loading" align="center">
      <p>Loading...</p>
    </div>
    <div v-else>
      <div v-for="(division, index) in divisionsTimetables" :key="index">
        <Timetable :divisionData="division" :divisionIndex="index" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import Timetable from "@/components/Timetable.vue";

const divisionsTimetables = ref([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const response = await fetch("/timetables.json"); // Path in public
    const data = await response.json();
    if (data.timetables) {
      divisionsTimetables.value = data.timetables; // Array of Days for each division
    }
  } catch (error) {
    console.error("Error loading timetables:", error);
  } finally {
    loading.value = false;
  }
});
</script>
