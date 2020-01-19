import React from 'react'
import { Popconfirm, message, Input, Select, Icon, Button, List, Row, Col } from 'antd'
import axios from 'axios'
import {NewServerDialog} from '../../components/dialog'
import GrpcServerComponent from './grpc'


export default class ServerApp extends React.Component {
  state = {current: undefined, loading: true, showNewDialog: false}
  servers = []
  protocol = 'grpc' // default new server protocol

  async componentDidMount() {
    await this.fetchServers()
    this.setState({current: this.servers[this.servers.length - 1], loading: false})
  }

  fetchServers = async () => {
    let resp
    try {
      resp = await axios.get('/api/v1/servers')
    } catch(err) {
      return message.error(err)
    }

    this.servers = resp.data
  }

  switchProtocol = (protocol) => {
    this.protocol = protocol
  }

  getServerComponent = (protocol) => {
    switch(protocol) {
      case 'grpc':
        return <GrpcServerComponent current={this.state.current} />
      case 'http':
        return <div>http server</div>
      case 'dubbo':
        return <div>dubbo server</div>
      default:
        return <Icon type='loading' />
    }
  }

  addServer = () => {
    this.setState({showNewDialog: true})
  }

  deleteServer = async (serverId) => {
    try {
      await axios.delete(`/api/v1/servers?id=${serverId}`)
      await this.fetchServers()
      this.setState({current: this.servers[this.servers.length - 1], loading: false})
    } catch (err) {
      return message.error(err.message.data)
    }
  }

  render() {
    const {current, loading, showNewDialog} = this.state

    return (
      <div>
        <Row>
          <Col md={4}>
            <Input.Group compact>
              <Select defaultValue={this.protocol} style={{width: '60%'}} onChange={this.switchProtocol}>
                <Select.Option value='grpc'>gRPC</Select.Option>
                <Select.Option value='http'>HTTP</Select.Option>
                <Select.Option value='dubbo'>Dubbo</Select.Option>
              </Select>
              <Button style={{width: '40%'}} onClick={this.addServer}>
                <Icon type='plus'/> New
              </Button>
            </Input.Group>
            <List
              size='small'
              dataSource={this.servers}
              bordered={false}
              loading={loading}
              renderItem={item => (
                <List.Item>
                  <Input.Group compact>
                    <Button
                      type='link'
                      style={(current!==undefined && current.id===item.id) ? {backgroundColor: '#1890FF', color: 'white', maxWidth: '80%'}: {}}
                      onClick={() => this.setState({current: item})}
                    >
                        # {item.id} {item.protocol} {item.addr}
                    </Button>
                    <Popconfirm
                      title='close this client?'
                      onConfirm={() => this.deleteClient(item.id)}
                    >
                      <Button
                        type='link'
                        style={{color: 'red', float: 'right'}}
                      >
                        <Icon type='delete'/>
                      </Button>
                    </Popconfirm>
                  </Input.Group>
                </List.Item>
              )}
            />
          </Col>
          <Col md={20} style={{paddingRight: 20, paddingLeft: 15}}>
            {
              current===undefined ?
                <Button style={{width: '100%'}} type='primary' onClick={this.addServer}>
                  <Icon type='plus'/> New Server
                </Button> :
                this.getServerComponent(current.protocol)
            }
          </Col>
        </Row>
        <NewServerDialog
          visible={showNewDialog}
          protocol={this.protocol}
          onClose={() => this.setState({showNewDialog: false})}
          onSubmit={async () => {
            await this.fetchServers()
            this.setState({
              current: this.servers[this.servers.length - 1],
              loading: false,
              showNewDialog: false})
          }}
        />
      </div>
    )
  }
}
