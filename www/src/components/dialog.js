import React from 'react'
import {Modal} from 'antd'
import GrpcClientForm from './grpcClientForm'
import HTTPClientForm from './httpClientForm'
import DubboClientForm from './dubboClientForm'


class NewClientDialog extends React.Component {
  getClientForm = protocol => {
    switch(protocol) {
      case 'grpc':
        return <GrpcClientForm onSubmit={this.props.onSubmit}/>
      case 'http':
        return <HTTPClientForm onSubmit={this.props.onSubmit}/>
      case 'dubbo':
        return <DubboClientForm onSubmit={this.props.onSubmit}/>
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
