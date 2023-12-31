import { useState, useEffect, useRef } from 'react'
import styled from 'styled-components'
import { capitalize } from 'lodash'
import { Button, Card, Space } from 'antd'
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'

interface LogState {
  index?: number;
  type: "success" | "error" | "warning" | "info";
  message: string[] | string;
}

const CartLog = styled(Card)`
  .ant-card-body {
    height: 500px;
  }
`

const LogTable = styled.div`
  table {
    border-spacing: 0;
    border-collapse: collapse;
    display: table;
    border-collapse: separate;
    box-sizing: border-box;
    text-indent: initial;
    border-spacing: 2px;
    border-color: grey;
  }
  td {
    padding: 0;
    display: table-cell;
    vertical-align: inherit;
  }
  .tab-size[data-tab-size="8"] {
    tab-size: 8;
  }
  .blob-num {
    position: relative;
    width: 1%;
    min-width: 50px;
    padding-right: 10px;
    padding-left: 10px;
    font-family: ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace;
    font-size: 12px;
    line-height: 20px;
    color: #6e7681;
    text-align: right;
    white-space: nowrap;
    vertical-align: top;
    cursor: pointer;
    -webkit-user-select: none;
    user-select: none;
  }
  .blob-num::before {
    content: attr(data-line-number);
  }
  .blob-code-inner {
    display: table-cell;
    overflow: visible;
    font-family: ui-monospace,SFMono-Regular,SF Mono,Menlo,Consolas,Liberation Mono,monospace;
    font-size: 12px;
    color: #c9d1d9;
    word-wrap: anywhere;
    white-space: pre;
  }
  .blob-code {
    position: relative;
    padding-right: 10px;
    padding-left: 10px;
    line-height: 20px;
    vertical-align: top;
  }
  .log {
    color: #c9d1d9;
  }
  .success {
    color: #28a745;
  }
  .error {
    color: #cb2431;
  }
  .info {
    color: #0366d6;
  }
  .debug {
    color: #6f42c1;
  }
  .warn {
    color: #f1e05a;
  }
`;

const Logs = () => {
  const [logs, setLogs] = useState<LogState[]>([]);
  const indexRef = useRef(0);
  const limit = 100;
  const addLog = (options: LogState) => {
    let {
      type,
      message
    } = options;
    if (logs.length > limit) {
      setLogs((prev) => [...prev.slice(1), {
        index: indexRef.current,
        type,
        message: typeof message == 'string' ? message : message.join(' ')
      }]);
    } else {
      setLogs((prev) => [...prev, {
        index: indexRef.current,
        type,
        message: typeof message == 'string' ? message : message.join(' ')
      }]);
    }
  };

  const clearLogs = () => {
    setLogs([]);
  };
  useEffect(() => {
    EventsOn('logger', (_, data: {
      level: "success" | "error" | "warning" | "info",
      message: string[] | string
    }) => {
      indexRef.current++;
      const message = typeof data.message === 'string' ? [data.message] : data.message.filter((item) => item != undefined);
      addLog({
        type: data.level,
        message: message.map((message) => typeof message === "string" ? message : JSON.stringify(message)).join(" - ")
      });
    })
    return () => {
      EventsOff('logger')
    }
  }, []);
  return (
    <>
      <CartLog
        title="Logs"
      >
        <Space style={{
          marginBottom: 10
        }}>
          <Button onClick={clearLogs}>Clear</Button>
        </Space>
        <LogTable style={{
          fontSize: 14,
          lineHeight: 1.5,
          color: "#c9d1d9",
          backgroundColor: "#1e1e1e",
          height: "calc(100% - 50px)",
          overflowY: "auto",
        }}>
          <table className="highlight tab-size" data-tab-size="8">
            <tbody>
              {
                logs?.map((log, index) => (
                  <tr key={index}>
                    <td className="blob-num">{log.index}</td>
                    <td className="blob-code blob-code-inner">
                      <span className={log.type}>{capitalize(log.type)}</span>: <span>{log.message}</span>
                    </td>
                  </tr>
                ))
              }
            </tbody>
          </table>
        </LogTable>
      </CartLog>
    </>
  );
};

export default Logs;