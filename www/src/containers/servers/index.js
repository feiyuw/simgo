import React from 'react'
import { Popconfirm, message, Input, Select, Icon, Button, List, Row, Col } from 'antd'
import axios from 'axios'
import {NewServerDialog} from '../../components/dialog'
import GrpcServerComponent from './grpc'
import { PROTOCOLS, GRPC, HTTP, DUBBO } from '../../constants'
import urls from '../../urls'


export default class ServerApp extends React.Component {
  state = {current: undefined, loading: true, showNewDialog: false}
  servers = []
  protocol = GRPC // default new server protocol

  async componentDidMount() {
    await this.fetchServers()
    this.setState({current: this.servers[0], loading: false})
  }

  fetchServers = async () => {
    let resp
    try {
      resp = await axios.get(urls.servers)
    } catch(err) {
      return message.error(err.response.data)
    }

    this.servers = resp.data
  }

  switchProtocol = (protocol) => {
    this.protocol = protocol
  }

  getServerComponent = (protocol) => {
    switch(protocol) {
      case GRPC:
        return <GrpcServerComponent current={this.state.current} />
      case HTTP:
        return <div>http server</div>
      case DUBBO:
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
      await axios.delete(urls.servers, {params: {id: serverId}})
      await this.fetchServers()
      this.setState({current: this.servers[0], loading: false})
    } catch (err) {
      return message.error(err.response.data)
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
                {
                  PROTOCOLS.map(p => (
                    <Select.Option key={p.value} value={p.value}>{p.name}</Select.Option>
                  ))
                }
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
                      style={(current!==undefined && current.name===item.name) ? {backgroundColor: '#1890FF', color: 'white', maxWidth: '80%'}: {}}
                      onClick={() => this.setState({current: item})}
                    >
                        # {item.id} {item.name} ({item.protocol}:{item.port})
                    </Button>
                    <Popconfirm
                      title='close this client?'
                      onConfirm={() => this.deleteServer(item.id)}
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
