<template>
  <el-menu
    :default-active="menuActive"
    class="el-menu-vertical-demo"
    background-color="#eeeeee"
    @select="handleSelected"
    :collapse="isCollapse"
    :unique-opened="true"
  >
    <SideBarItem v-for="(menu,idx) in menuList" :key="idx" :item="menu"/>
    <!-- <el-menu-item index="/v1/home">
      <i class="el-icon-s-home"></i>
      <span>首页</span>
    </el-menu-item>
    <el-submenu index="/ca">
      <template slot="title">
        <i class="el-icon-document"></i>
        <span>证书管理</span>
      </template>
      <el-menu-item index="/v1/newca/server">
        <i class="el-icon-menu"></i>
        <span>服务器证书</span>
      </el-menu-item>
      <el-menu-item index="/v1/newca/personal">
        <i class="el-icon-menu"></i>
        <span>个人证书</span>
      </el-menu-item>
      <el-menu-item index="/v1/newca/sign">
        <i class="el-icon-menu"></i>
        <span>代码签名证书</span>
      </el-menu-item>
    </el-submenu>
    <el-submenu index="3">
      <template slot="title">
        <i class="el-icon-document-delete"></i>
        <span>获取吊销证书信息</span>
      </template>
      <el-menu-item index="3-1">
        <i class="el-icon-menu"></i>
        <span>用户证书列表</span>
      </el-menu-item>
      <el-menu-item index="3-2">
        <i class="el-icon-menu"></i>
        <span>获取用户证书</span>
      </el-menu-item>
      <el-menu-item index="3-3">
        <i class="el-icon-menu"></i>
        <span>验证证书信息</span>
      </el-menu-item>
    </el-submenu>
    <el-submenu index="/v1/setting">
      <template slot="title">
        <i class="el-icon-setting"></i>
        <span>配置</span>
      </template>
      <el-menu-item index="/v1/setting/user">
        <i class="el-icon-menu"></i>
        <span>用户权限</span>
      </el-menu-item>
    </el-submenu>
    <el-menu-item index="/v1/help">
      <i class="el-icon-help"></i>
      <span>使用帮助</span>
    </el-menu-item>
    <el-menu-item index="/v1/system">
      <i class="el-icon-notebook-1"></i>
      <span>系统信息</span>
    </el-menu-item> -->
  </el-menu>
</template>

<script>
import SideBarItem from './sidebaritem/index'

export default {
  components: {
    SideBarItem
  },
  data() {
    return {
      menuActive: "",
      isCollapse: false,
      menuList: []
    };
  },

  created() {
    this.menuActive = this.getUrlPath()
    this.getMenuList()
  },

  watch: {
    $route(to, from) {
      let path = to.path;
      this.menuActive = path;
    }
  },

  methods: {
    getUrlPath() {
      let url = window.location.href.split("//")
      let path = url[1]
      return path.substring(path.indexOf("/"))
    },
    handleSelected(index, indexPath) {
      let oldIndex = this.getUrlPath()
      if (index === oldIndex) return
      this.$router.push(index);
    },
    async getMenuList() {
      try {
        let res = await this.$api.menu.GetMenuList()
        this.menuList = res.items
      } catch (error) {
        console.log(error)
      }
    }
  }
};
</script>

<style lang="scss" scoped></style>
