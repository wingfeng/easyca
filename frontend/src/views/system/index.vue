<template>
  <div>
    <el-card shadow="never">
      <div slot="header" class="clearfix">
        系统信息
      </div>
      <div class="list">
        <div>
          <el-row>
            <el-col>
              <span>系统名称：</span>
              <span>{{ systemInfo.AppName }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>系统版本：</span>
              <span>{{ systemInfo.AppVersion }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>Go版本：</span>
              <span>{{ systemInfo.GoVersion }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>Git修订：</span>
              <span>{{ systemInfo.GitRevision }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>Tag分支: </span>
              <span>{{ systemInfo.GitBranch }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>编译版本：</span>
              <span>{{ systemInfo.BuildVersion }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>编译时间：</span>
              <span>{{ systemInfo.BuildTime }}</span>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: "System",
  components: {},
  data() {
    return {
      systemInfo: {}
    };
  },
  created() {
    this.getSystemVersion();
  },
  methods: {
    async getSystemVersion() {
      try {
        let res = await this.$api.system.SystemVersion();
        if (res.BuildTime) {
          res.BuildTime = this.$moment(res.BuildTime).format(
            "YYYY-MM-DD HH:mm:ss"
          );
        }
        this.systemInfo = res;
      } catch (error) {
        console.log(error);
      }
    }
  }
};
</script>

<style lang="less" scoped>
@topOffset: 50%;
@leftOffset: 50%;

.list {
  position: relative;
  height: 70vh;
  line-height: 40px;

  div > .el-row > .el-col > span:first-child {
    width:80px;
    display: inline-block;
    text-align: right;
    margin-right: 10px;
  }
}
</style>
