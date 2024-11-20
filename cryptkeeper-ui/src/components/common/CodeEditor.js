import React, { useState, useEffect } from 'react';
import CodeMirror from "@uiw/react-codemirror";
import { vscodeDark } from "@uiw/codemirror-theme-vscode";
import { json } from '@codemirror/lang-json';

import { oneDark } from '@codemirror/theme-one-dark';


function CodeEditor({ code, onChange, height, extensions }) {


    return (
        <div className="code_editor">
            <CodeMirror
                className='bg-dark p-2 rounded-2'
                value={code}
                extensions={extensions ? extensions : [json()]}
                theme={oneDark}
                onChange={onChange}
                height={ height ? height: "400px" }
                options={{
                    inlineSuggest: true,
                    fontSize: "12px",
                    formatOnType: true,
                    indentUnit: 2,
                    smartIndent: true,
                    autoClosingBrackets: true,
                    minimap: { scale: 10 },
                }}

            />
        </div>
    );
}

export default CodeEditor;
