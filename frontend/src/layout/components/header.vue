<template>
  <div class="header">
    <div class="logo">
      <h1 @click="toHome">EASY CA</h1>
      <i class="el-icon-s-unfold menu-icon" @click="handleMenu"></i>
    </div>
    <div>
      <el-dropdown @command="handleCommand">
        <span>{{ user.userInfo.name }}<i class="el-icon-caret-bottom"></i> </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="logout">退出</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: "Header",
  created() {
    this.getUserInfo()
  },
  data() {
    return {
      userInfo: {}
    }
  },
  computed: {
    ...mapGetters([
      'user'
    ])
  },
  methods: {
    toHome() {
      window.location.href = '/'
    },
    async getUserInfo() {
      try {
        let res = await this.$api.user.UserInfo()
        if (!res.name) res.name = "游客"
        this.$store.dispatch("user/setUser", res)
      }
      catch (err) {
        console.log(err)
      }
    },
    handleCommand(type) {
      if (type === 'logout') {
        window.location.href = window.CONFIG.logoutURL
      }
    },
    handleMenu() {
      this.$emit("handleMenu")
    }
  }
};
</script>

<style scoped lang="less">
.headerColor {
  color: #fff
}

.header {
  background: #000;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 45px;
  .headerColor();

  h1 {
    font-size: 18px;
    font-weight: normal;
    cursor: pointer;
  }

  .el-dropdown {
    .headerColor();
  }
}

.logo {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.menu-icon {
  font-size: 23px;
  margin-left: 40px;
  cursor: pointer;
  @media only screen and (min-width: 400px){
    display: block;
  }

  @media only screen and (min-width: 800px){
    display: none;
  }
}
</style>
