<template>
  <div>
    <!-- first menu -->
    <el-menu-item v-if="item.subs === null" :index="item.url">
      <i :class="item.class"></i>
      <span>{{item.name}}</span>
    </el-menu-item>
    <!-- second menu -->
    <el-submenu v-else :index="item.url">
      <template slot="title">
        <i :class="item.class"></i>
        <span>{{item.name}}</span>
      </template>
      <!-- children menu -->
      <el-menu-item v-for="(children,idx) in item.subs" :key="idx" :index="children.url">
        <i class="el-icon-menu"></i>
        <span>{{children.name}}</span>
      </el-menu-item>
    </el-submenu>
  </div>
</template>

<script>
export default {
  name: "SideBarItem",
  props: {
    item: {
      type: Object,
      required: true
    }
  }
};
</script>

<style lang="less" scoped></style>
