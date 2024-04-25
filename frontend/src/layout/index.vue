<template>
  <div>
    <Header @handleMenu="handleMenu"/>
    <div class="main">
      <!-- adapt to the mobile phone -->
      <el-collapse-transition>
        <div class="menu" v-if="menuFlag">
          <Sidebar />
        </div>
      </el-collapse-transition>
      <!-- adapt to the mobile phone -->
      <el-row>
        <el-col :xs="0" :sm="4" :md="4" :lg="4" :xl="3">
          <Sidebar />
        </el-col>
        <el-col
          class="content-container"
          :xs="24"
          :sm="20"
          :md="20"
          :lg="20"
          :xl="21"
        >
          <router-view></router-view>
        </el-col>
      </el-row>
    </div>
    <Footer />
  </div>
</template>

<script>
import Header from "./components/header.vue";
import Sidebar from "./components/sidebar.vue";
import Footer from "./components/footer.vue";

export default {
  components: {
    Header,
    Sidebar,
    Footer
  },
  data() {
    return {
      menuFlag: false
    }
  },
  methods: {
    handleMenu() {
      this.menuFlag = !this.menuFlag
    }
  }
};
</script>

<style lang="less" scoped>
.main {
  position: relative;
}
.content-container {
  margin-top: 30px;
  padding: 0 20px;
}
.menu {
  position: absolute;
  width: 50%;
  height: 100vh;
  left: 0;
  z-index: 99;
}
</style>
