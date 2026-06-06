import { useEffect, useState, useMemo, useCallback } from "react";
import ExampleTheme from "./themes/ExampleTheme";
import {$getRoot, $getSelection, $createParagraphNode} from 'lexical';
import { LexicalComposer } from "@lexical/react/LexicalComposer";
import { RichTextPlugin } from "@lexical/react/LexicalRichTextPlugin";
import { ContentEditable } from "@lexical/react/LexicalContentEditable";
import { HistoryPlugin } from "@lexical/react/LexicalHistoryPlugin";
import { AutoFocusPlugin } from "@lexical/react/LexicalAutoFocusPlugin";
import LexicalErrorBoundary from "@lexical/react/LexicalErrorBoundary";
import { $generateHtmlFromNodes, $generateNodesFromDOM } from '@lexical/html';
import TreeViewPlugin from "./plugins/TreeViewPlugin";
import ToolbarPlugin from "./plugins/ToolbarPlugin";
import { HeadingNode, QuoteNode } from "@lexical/rich-text";
import { TableCellNode, TableNode, TableRowNode } from "@lexical/table";
import { ListItemNode, ListNode } from "@lexical/list";
import { CodeHighlightNode, CodeNode } from "@lexical/code";
import { AutoLinkNode, LinkNode } from "@lexical/link";
import { LinkPlugin } from "@lexical/react/LexicalLinkPlugin";
import { ListPlugin } from "@lexical/react/LexicalListPlugin";
import { MarkdownShortcutPlugin } from "@lexical/react/LexicalMarkdownShortcutPlugin";
import { TRANSFORMERS } from "@lexical/markdown";
import {useLexicalComposerContext} from '@lexical/react/LexicalComposerContext';
import { ClearEditorPlugin } from '@lexical/react/LexicalClearEditorPlugin';
import ListMaxIndentLevelPlugin from "./plugins/ListMaxIndentLevelPlugin";
import CodeHighlightPlugin from "./plugins/CodeHighlightPlugin";
import AutoLinkPlugin from "./plugins/AutoLinkPlugin";
import ActionsPlugin from "../editor/plugins/ActionsPlugin"
import { debounceLatest } from "../../amper/Instruments";
import Convenience from "../../help/Convenience";
import "./ChatEditor.css";

function Placeholder() {
  return <div className="editor-placeholder">Enter your new message here...</div>;
}

function OnChangePlugin({ onChange }) {
  const [editor] = useLexicalComposerContext();
  useEffect(() => {
    return editor.registerUpdateListener(({editorState}) => {
      onChange(editor, editorState);
    });
  }, [editor, onChange]);
}

export default function ChatEditor({content, send, onAttachmentClicked, showSaveCancel, save, cancel, sendEnabled}) {

  const debounceUpdateContent = (arg0, content) => {
    setState({
      ...state,
      content,
    });
  }

  const [state, setState] = useState({
    content: '',
    debounceUpdateContent: debounceLatest(debounceUpdateContent, 500),
  });

  const onChange = (editor, editorState) => {
    editor.update(() => {
      const htmlString = $generateHtmlFromNodes(editor, null);
      setState({
        ...state,
        content: htmlString,
      });
   });
  };

  const initEditor = useCallback((editor) => {
    if (content) {
      // In the browser you can use the native DOMParser API to parse the HTML string.
      const parser = new DOMParser();
      const dom = parser.parseFromString(content, "text/html");

      // Once you have the DOM instance it's easy to generate LexicalNodes.
      const nodes = $generateNodesFromDOM(editor, dom);
      // Select the root
      //const para = $createParagraphNode();
      //para.append(...nodes);
      $getRoot().append(...nodes);
    }
  }, []);

  const initialConfig = useMemo(
    () => ({
      editorState: initEditor,
        // The editor theme
      theme: ExampleTheme,
      // Handling of errors during update
      onError(error) {
        throw error;
      },
      // Any custom nodes go here
      nodes: [
        HeadingNode,
        ListNode,
        ListItemNode,
        QuoteNode,
        CodeNode,
        CodeHighlightNode,
        TableNode,
        TableCellNode,
        TableRowNode,
        AutoLinkNode,
        LinkNode
      ]
    }),
    []
  );
  
  const extractContent = (html) => {
    var span = document.createElement('span');
    span.innerHTML = html;
    return span.textContent || span.innerText;
  };

  const sendInternal = () => {
    send(state.content);
    setState({
      ...state,
      content: '',
    });
  };

  return (
    <LexicalComposer initialConfig={initialConfig}>
      <div className="editor-container">
        <ToolbarPlugin style={{overflow: 'auto'}} onAttachmentClicked={onAttachmentClicked} />
        <div className="editor-inner">
          <RichTextPlugin
            contentEditable={<ContentEditable className="editor-input" />}
            placeholder={<Placeholder />}
            ErrorBoundary={LexicalErrorBoundary}
          />
          <OnChangePlugin onChange={onChange}/>
          <HistoryPlugin />
          <AutoFocusPlugin />
          <CodeHighlightPlugin />
          <ClearEditorPlugin />
          <ListPlugin />
          <LinkPlugin />
          <AutoLinkPlugin />
          <ListMaxIndentLevelPlugin maxDepth={7} />
          <MarkdownShortcutPlugin transformers={TRANSFORMERS} />
        </div>
        <ActionsPlugin 
          showSaveCancel={showSaveCancel} 
          send={sendInternal} 
          sendEnabled={(Convenience.hasValue(extractContent(state.content)) || sendEnabled)}
          save={()=>{save(state.content)}} 
          cancel={cancel}/>
      </div>
    </LexicalComposer>
  );
}
