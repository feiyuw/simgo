import React from 'react'
import {Modal} from 'antd'
import GrpcClientForm from './grpcClientForm'
import HTTPClientForm from './httpClientForm'
import DubboClientForm from './dubboClientForm'
import GrpcServerForm from './grpcServerForm'
import HTTPServerForm from './httpServerForm'
import DubboServerForm from './dubboServerForm'
import GrpcHandlerForm from './grpcHandlerForm'


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


class NewServerDialog extends React.Component {
  getServerForm = protocol => {
    switch(protocol) {
      case 'grpc':
        return <GrpcServerForm onSubmit={this.props.onSubmit}/>
      case 'http':
        return <HTTPServerForm onSubmit={this.props.onSubmit}/>
      case 'dubbo':
        return <DubboServerForm onSubmit={this.props.onSubmit}/>
      default:
        return <div>Unsupported!</div>
    }
  }

  render() {
    return (
      <Modal
        title={`New ${this.props.protocol} server`}
        visible={this.props.visible}
        onCancel={this.props.onClose}
        footer={null}
        width='50%'
        destroyOnClose={true}
        confirmLoading={true}
      >
        {this.getServerForm(this.props.protocol)}
      </Modal>
    )
  }
}


class NewGrpcMethodHandlerDialog extends React.Component {
  render() {
    return (
      <Modal
        title={`New Grpc Method Handler`}
        visible={this.props.visible}
        onCancel={this.props.onClose}
        footer={null}
        width='50%'
        destroyOnClose={true}
        confirmLoading={true}
      >
      <GrpcHandlerForm
        server={this.props.server}
        onSubmit={this.props.onSubmit}
      />
      </Modal>
    )
  }
}

export { NewClientDialog, NewServerDialog, NewGrpcMethodHandlerDialog }
