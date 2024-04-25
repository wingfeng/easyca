<template>
  <div>
    <el-card shadow="never">
      <div slot="header" class="clearfix">
        个人证书信息
      </div>
      <div class="list" v-loading="loading" element-loading-text="创建中">
        <div v-if="haveCA">
          <el-row>
            <el-col>
              <span>版本：</span>
              <span>V{{ caInfo.Version }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>序列号：</span>
              <span>{{ caInfo.SerialNumber }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <span>个人地址：</span>
              <span>{{ caInfo.Subject.CommonName }}</span>
            </el-col>
          </el-row>
          <el-row>
            <el-col
              ><el-button
                type="primary"
                size="mini"
                @click="download"
                icon="el-icon-download"
                >下载个人证书</el-button
              ></el-col
            >
          </el-row>
        </div>
        <div v-else class="ca-create">
          <el-link type="info" @click="createCA"
            >暂无个人证书，点击进行创建</el-link
          >
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { zipDownload } from "@/utils/download";

export default {
  name: "Personal",
  components: {},
  data() {
    return {
      haveCA: false,
      caInfo: null,
      loading: false
    };
  },
  created() {
    this.getPersonal();
  },
  methods: {
    createServer() {
      this.$refs.form.show();
    },
    async getPersonal() {
      try {
        let res = await this.$api.personal.GetPersonal();
        this.haveCA = true;
        this.caInfo = res;
      } catch (error) {
        this.haveCA = false;
        if (error.data.errInfo) {
          this.$message["warning"](error.data.errInfo);
        } else {
          this.$message["error"](error.message);
        }
      }
    },
    async createCA() {
      this.loading = true;
      try {
        let res = await this.$api.personal.CreatePersonal();
        if (res) {
          this.loading = false;
          this.$message["success"]("创建成功");
          this.getPersonal();
        }
      } catch (error) {
        this.loading = false;
        if (error.data.errInfo) {
          this.$message["warning"](error.data.errInfo);
        } else {
          this.$message["error"](error.message);
        }
      }
    },
    async download() {
      try {
        let res = await this.$api.personal.DownloadCA();
        zipDownload(res);
      } catch (error) {
        if (error.data.errInfo) {
          this.$message["warning"](error.data.errInfo);
        } else {
          this.$message["error"](error.message);
        }
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

  .ca-create {
    position: absolute;
    top: @topOffset;
    left: @leftOffset;
    transform: translate(-@leftOffset, -@topOffset);
  }

  div > .el-row > .el-col > span:first-child {
    width: 80px;
    display: inline-block;
    text-align: right;
    margin-right: 10px;
  }
}
</style>
