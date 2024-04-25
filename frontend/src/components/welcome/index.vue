<template>
  <div>
    <p>{{ welcomeContent() }}</p>
  </div>
</template>

<script>
import timeFormat from "@/utils/time";
export default {
  name: "Welcome",
  props: {
    name: {
      type: String,
      default: "visitor"
    }
  },
  methods: {
    welcomeContent() {
      let date = new Date();
      let hour = date.getHours();
      let time = timeFormat(date);
      let name = this.name;
      if (hour >=0 && hour < 8) return `早上好：${name}！当前时间：${time}`;
      else if (hour >= 8 && hour < 11) return `上午好：${name}！当前时间：${time}`;
      else if (hour >= 11 && hour < 13) return `中午好：${name}！当前时间：${time}`;
      else if (hour >= 13 && hour < 17) return `下午好：${name}！当前时间：${time}`;
      else if (hour >= 17 && hour < 19) return `傍晚好：${name}！当前时间：${time}`;
      else return `晚上好：${name}！当前时间：${time}`;
    }
  }
};
</script>

<style lang="scss" scoped></style>
