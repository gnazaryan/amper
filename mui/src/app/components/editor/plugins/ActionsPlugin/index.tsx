import {useLexicalComposerContext} from '@lexical/react/LexicalComposerContext';
import Button from '@mui/material/Button';
import SendIcon from '@mui/icons-material/Send';
import {Dispatch} from 'react';
import SaveIcon from '@mui/icons-material/Save';
import CancelIcon from '@mui/icons-material/Cancel';
import { CLEAR_EDITOR_COMMAND } from 'lexical';

export default function ActionsPlugin({
  send,
  sendEnabled,
  save,
  cancel,
  showSaveCancel,
}: {
    send: Dispatch<void>;
    sendEnabled: boolean;
    save: Dispatch<void>;
    cancel: Dispatch<void>;
    showSaveCancel: boolean;
}): JSX.Element {
  const [editor] = useLexicalComposerContext();

    const sendInternal = () => {
      editor.dispatchCommand(CLEAR_EDITOR_COMMAND, undefined);
      editor.focus();
      send();
    };

    const getView = () => {
        if (showSaveCancel === true) {
            return [<Button sx={{mr: 1}} variant="contained" endIcon={<CancelIcon />} onClick={()=>{cancel()}}>
                Cancel
            </Button>,<Button variant="contained" endIcon={<SaveIcon />} onClick={()=>{save()}}>
                Save
            </Button>];
        } else {
            return <Button variant="contained" endIcon={<SendIcon />} onClick={sendInternal} disabled={!sendEnabled}>
                Send
            </Button>;
        }
    };

  return (
    <div className="actions">
      {getView()}
    </div>
  );
}