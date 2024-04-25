<template>
  <div>
    <el-card shadow="never">
      <div slot="header" class="clearfix">
        使用帮助
      </div>
      <div v-html="help"></div>
    </el-card>
  </div>
</template>

<script>
import help from "./HELP.md";
export default {
  data() {
    return {
      help
    };
  }
};
</script>

<style lang="scss" scoped></style>
