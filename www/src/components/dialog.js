import React from 'react'
import axios from 'axios'
import {Modal, message} from 'antd'
import GrpcClientForm from './grpcClientForm'
import HTTPClientForm from './httpClientForm'
import DubboClientForm from './dubboClientForm'


class NewClientDialog extends React.Component {
  handleSubmit = client => {
    // new
    axios.post('/api/v1/clients', {})
      .then(resp => {
        message.success(`Client # ${resp.data.id} added!`)
        this.props.onSubmit(resp.data)
      })
      .catch(err => {
        message.error(`add client error，${err.response.data}`)
      })
  }

  getClientForm = protocol => {
    switch(protocol) {
      case 'grpc':
        return <GrpcClientForm />
      case 'http':
        return <HTTPClientForm />
      case 'dubbo':
        return <DubboClientForm />
      default:
        return <div>Unsupported!</div>
    }
  }

  render() {
    return (
      <Modal
        title={`New ${this.props.protocol} client`}
        visible={this.props.visible}
        onCancel={this.props.onClose}
        footer={null}
        width='50%'
        destroyOnClose={true}
        confirmLoading={true}
      >
        {this.getClientForm(this.props.protocol)}
      </Modal>
    )
  }
}

export { NewClientDialog }
