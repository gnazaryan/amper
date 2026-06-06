import Popover from '@mui/material/Popover';
import emoticons from './Emoticons';
import Typography from '@mui/material/Typography';
import './EmotionPicker.css'

export default function EmotionPicker({el, onClose, onSelect}) {

    const getEmoticons = () => {
        const result = [];
        for (const [category, categoryEmoticons] of Object.entries(emoticons)) {
            const categoryResult = [];
            for (let i = 0; i < categoryEmoticons.length; i++) {
                categoryResult.push(
                    <span title={categoryEmoticons[i].shortName} className="emoticon" onClick={() => {onSelect(categoryEmoticons[i].emoticon)}}>{categoryEmoticons[i].emoticon}</span>
                );
            }
            result.push(
                <div id={category}>
                    <Typography variant="caption" display="block" gutterBottom sx={{ml: 1}}>
                        #{category}
                    </Typography>
                    {categoryResult}
                </div>
            );
        }

        return result;
    };

    return <Popover
        open={Boolean(el)}
        anchorEl={el}
        onClose={onClose}
        anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
        }}
        sx={{
            maxHeight: 250,
            maxWidth: 370,
        }}
    >
        {getEmoticons()}
    </Popover>;
}