import Convenience from "../../../help/Convenience";
import { formatDate } from "../../../amper/Instruments";

export const FlagSeen = "\\Seen"
export const FlagAnswered = "\\Answered"
export const FlagFlagged  = "\\Flagged"
export const FlagDeleted  = "\\Deleted"
export const FlagDraft    = "\\Draft"

// Widely used flags
export const FlagForwarded = "$Forwarded"
export const FlagMDNSent   = "$MDNSent" // Message Disposition Notification sent
export const FlagJunk      = "$Junk"
export const FlagNotJunk   = "$NotJunk"
export const FlagPhishing  = "$Phishing"
export const FlagImportant = "$Important" // RFC 8457

// Permanent flags
export const FlagWildcard = "\\*"

export const getAddress = (addresses) => {
    const result = [];
    if (addresses) {
        for (let i = 0; i < addresses.length; i++) {
            result.push(addresses[i].mailbox + '@' + addresses[i].host);
        }
    }
    return result.join(",");;
};

export const getFrom = (email) => {
    if (email.envelope && email.envelope.from && email.envelope.from.length > 0) {
        const result = [];
        for (let i = 0; i < email.envelope.from.length; i++) {
            if (Convenience.hasValue(email.envelope.from[i].name)) {
                result.push(email.envelope.from[i].name);
            } else if (Convenience.hasValue(email.envelope.from[i].mailbox) && Convenience.hasValue(email.envelope.from[i].host)) {
                result.push(email.envelope.from[i].mailbox + '@' + email.envelope.from[i].host);
            }
        }
        return result.join(', ');
    }
    return 'Unknown';
}

export const getFromName = (email) => {
    if (email.envelope && email.envelope.from && email.envelope.from.length > 0) {
        const result = [];
        for (let i = 0; i < email.envelope.from.length; i++) {
            if (Convenience.hasValue(email.envelope.from[i].name)) {
                result.push(email.envelope.from[i].name);
            }
        }
        return result.join(', ');
    }
    return 'Unknown';
}

export const getFromEmail = (email) => {
    if (email.envelope && email.envelope.from && email.envelope.from.length > 0) {
        const result = [];
        for (let i = 0; i < email.envelope.from.length; i++) {
            if (Convenience.hasValue(email.envelope.from[i].mailbox) && Convenience.hasValue(email.envelope.from[i].host)) {
                result.push(email.envelope.from[i].mailbox + '@' + email.envelope.from[i].host);
            }
        }
        return result.join(', ');
    }
    return 'Unknown';
}

export const getDate = (email) => {
    if (email.envelope) {
        return formatDate(email.envelope.date);
    }
};

export const getSubject = (email) => {
    if (email.envelope) {
        return email.envelope.subject;
    }
    return "Unknown"
};

export const isSeen = (email) => {
    if (email && email.flags && email.flags.indexOf(FlagSeen) > -1) {
        return true;
    }
    return false;
};