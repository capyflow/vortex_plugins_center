<template>
  <div class="health-indicator" :class="status">
    <div class="pulse"></div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  status?: 'healthy' | 'unhealthy' | string
}>()
</script>

<style scoped>
.health-indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  position: relative;
}

.health-indicator.healthy {
  background: #2ed573;
}

.health-indicator.unhealthy {
  background: #ff4757;
}

.pulse {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100%;
  height: 100%;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.healthy .pulse {
  background: rgba(46, 213, 115, 0.4);
}

.unhealthy .pulse {
  background: rgba(255, 71, 87, 0.4);
}

@keyframes pulse {
  0% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) scale(2.5);
    opacity: 0;
  }
}
</style>
