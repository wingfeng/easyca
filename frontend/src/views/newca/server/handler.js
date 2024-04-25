export const handler = {
  methods: {
    create() {
      this.model({
        title: '新建banner',
        callback: this.createData
      })
    },
    edit(data) {
      this.model({
        title: '修改banner',
        data: Object.assign({}, data),
        callback: this.updateData
      })
    },
    deleteBatch() {
      const selectedData = this.getSelected()
      if (selectedData.length !== 0) {
        this.deleteData(
          selectedData.map(row => row.id),
          '确认批量删除选中数据吗？'
        )
      } else {
        this.$message({
          message: '请先选择数据',
          type: 'error'
        })
      }
    },
    createData(valid, data, model) {
      if (valid) {
        this.$api.school_banner.Add(data)
          .then(res => {
            if (res.code === 200) {
              this.$message['success']('添加成功')
              this.refresh()
              model.close()
            } else {
              return Promise.reject()
            }
          }).catch(() => {
            this.$message['error']('添加失败')
          })
      }
    },
    updateData(valid, data, model) {
      if (valid) {
        this.$api.school_banner.Update(data)
          .then(res => {
            if (res.code === 200) {
              this.$message['success']('修改成功')
              this.refresh()
              model.close()
            } else {
              return Promise.reject()
            }
          }).catch((err) => {
            console.log('err', err)
            this.$message['error']('修改失败')
          })
      }
    },
    deleteData(ids, tip = '确认删除？') {
      this.$confirm(tip, '提示', {
        type: 'warning'
      }).then(() => {
        const data = ids.map(id => ({ id }))
        this.$api.school_banner.Delete(data)
          .then(res => {
            if (res.code === 200) {
              this.$message['success']('删除成功')
              this.refresh()
            } else {
              return Promise.reject()
            }
          }).catch(() => {
            this.$message['error']('删除失败')
          })
      })
    }
  }
}
