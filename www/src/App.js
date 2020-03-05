import React from 'react'
import { Layout, Menu } from 'antd'
import "antd/dist/antd.css"
import { HashRouter, Switch, Route, Link } from 'react-router-dom'
import ServerApp from './containers/servers'
import ClientApp from './containers/clients'

const { Header, Content, Footer } = Layout


function App() {
  return (
    <HashRouter>
      <Layout>
        <Header style={{ position: 'fixed', zIndex: 1, width: '100%' }}>
          <div className="logo" />
          <Menu
            theme="dark"
            mode="horizontal"
            defaultSelectedKeys={[window.location.hash.substr(1) === "/" ? "/client" : window.location.hash.substr(1)]}
            style={{ lineHeight: '64px' }}
          >
            <Menu.Item key="/client"><Link to='/client'>clients</Link></Menu.Item>
            <Menu.Item key="/server"><Link to='/server'>servers</Link></Menu.Item>
          </Menu>
        </Header>
        <Content style={{ padding: '0 50px', marginTop: 64 }}>
          <div style={{ background: '#fff', padding: 24, minHeight: 480 }}>
            <Switch>
              <Route key='/client' path='/' exact component={ClientApp}/>
              <Route key='/client' path='/client' exact component={ClientApp}/>
              <Route key='/server' path='/server' exact component={ServerApp}/>
            </Switch>
          </div>
        </Content>
        <Footer style={{ textAlign: 'center' }}>Simgo ©2020 Created by feiyuw</Footer>
      </Layout>
    </HashRouter>
  )
}

export default App
